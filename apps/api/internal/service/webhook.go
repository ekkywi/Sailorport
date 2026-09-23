package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ekkywi/sailorport/apps/api/internal/model"
)

type webhookCatalog interface {
	ListAll(ctx context.Context) ([]model.Service, error)
}

type webhookDeployer interface {
	Create(ctx context.Context, serviceID string, req model.CreateDeploymentRequest, actorID, role string) (model.Deployment, error)
}

type webhookDeliveries interface {
	Exists(ctx context.Context, deliveryID string) (bool, error)
	Insert(ctx context.Context, deliveryID, serviceID, deploymentID string) (model.WebhookDelivery, error)
}

type Webhook struct {
	catalog     webhookCatalog
	deployments webhookDeployer
	deliveries  webhookDeliveries
}

func NewWebhook(catalog webhookCatalog, deployments webhookDeployer, deliveries webhookDeliveries) *Webhook {
	return &Webhook{
		catalog:     catalog,
		deployments: deployments,
		deliveries:  deliveries,
	}
}

// HandleGitHub verifies a GitHub push webhook and may create a deployment when auto-deploy is on.
//
// Signature is verified before catalog-dependent outcomes are returned. Duplicate
// X-GitHub-Delivery ids are ignored after a successful prior Create+Insert.
func (w *Webhook) HandleGitHub(
	ctx context.Context,
	event string,
	deliveryID string,
	signatureHeader string,
	body []byte,
) (model.WebhookAck, error) {
	event = strings.TrimSpace(event)
	deliveryID = strings.TrimSpace(deliveryID)
	if event == "" {
		return model.WebhookAck{}, fmt.Errorf("%w: missing X-GitHub-Event", ErrInvalid)
	}

	ack := model.WebhookAck{Received: true, Event: event}

	if event != "push" {
		ack.Ignored = true
		ack.Reason = "only push events are handled"
		return ack, nil
	}

	var payload model.GitHubPushPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		return model.WebhookAck{}, fmt.Errorf("%w: invalid JSON body", ErrInvalid)
	}

	branch := branchFromRef(payload.Ref)
	if branch == "" {
		ack.Ignored = true
		ack.Reason = "not a branch ref (tag or empty)"
		return ack, nil
	}

	matches, err := w.findServicesByCloneURL(ctx, payload.Repository.CloneURL)
	if err != nil {
		return model.WebhookAck{}, err
	}

	authed := filterVerifiedServices(matches, body, signatureHeader)
	if len(authed) == 0 {
		return model.WebhookAck{}, fmt.Errorf("%w: invalid signature", ErrUnauthorized)
	}

	ack.Repo = payload.Repository.FullName
	ack.CloneURL = payload.Repository.CloneURL
	ack.Branch = branch
	ack.CommitSHA = payload.After
	ack.Pusher = payload.Pusher.Name

	eligible := filterAutoDeployServices(authed, branch)
	if len(eligible) == 0 {
		ack.Ignored = true
		ack.Reason = "no service with auto-deploy enabled for this branch"
		return ack, nil
	}

	target := eligible[0]
	env := strings.TrimSpace(target.AutoDeployEnvironment)
	if env == "" {
		env = "staging"
	}

	if w.deployments == nil {
		return model.WebhookAck{}, fmt.Errorf("webhook deployer not configured")
	}

	if deliveryID != "" && w.deliveries != nil {
		exists, err := w.deliveries.Exists(ctx, deliveryID)
		if err != nil {
			return model.WebhookAck{}, err
		}
		if exists {
			ack.Ignored = true
			ack.Reason = "duplicate delivery"
			return ack, nil
		}
	}

	dep, err := w.deployments.Create(ctx, target.ID, model.CreateDeploymentRequest{
		Environment: env,
	}, "", "")
	if err != nil {
		return model.WebhookAck{}, err
	}

	if deliveryID != "" && w.deliveries != nil {
		if _, err := w.deliveries.Insert(ctx, deliveryID, target.ID, dep.ID); err != nil {
			// Rare race: another request already recorded this id — OK.
			if ok, exErr := w.deliveries.Exists(ctx, deliveryID); exErr != nil || !ok {
				return model.WebhookAck{}, fmt.Errorf("record webhook delivery: %w", err)
			}
		}
	}

	ack.ServiceID = target.ID
	ack.DeploymentID = dep.ID
	ack.Environment = env
	return ack, nil
}

func branchFromRef(ref string) string {
	const prefix = "refs/heads/"
	if !strings.HasPrefix(ref, prefix) {
		return ""
	}
	return strings.TrimPrefix(ref, prefix)
}

func normalizeRepoURL(u string) string {
	u = strings.TrimSpace(strings.ToLower(u))
	u = strings.TrimSuffix(u, "/")
	u = strings.TrimSuffix(u, ".git")
	return u
}

func (w *Webhook) findServicesByCloneURL(ctx context.Context, cloneURL string) ([]model.Service, error) {
	if w.catalog == nil {
		return nil, fmt.Errorf("webhook catalog not configured")
	}
	all, err := w.catalog.ListAll(ctx)
	if err != nil {
		return nil, err
	}
	want := normalizeRepoURL(cloneURL)
	var out []model.Service
	for _, svc := range all {
		if svc.SourceType != "git" {
			continue
		}
		if normalizeRepoURL(svc.RepoURL) == want {
			out = append(out, svc)
		}
	}
	return out, nil
}

// filterVerifiedServices menyisakan service yang secret-nya sendiri cocok dengan
// signature body ini. Deploy hanya boleh berjalan untuk service pemegang secret —
// dua service yang menunjuk repo sama tidak boleh saling memicu deploy.
func filterVerifiedServices(services []model.Service, body []byte, signatureHeader string) []model.Service {
	var out []model.Service
	for _, svc := range services {
		if err := verifyGitHubSignature(svc.WebhookSecret, body, signatureHeader); err != nil {
			continue
		}
		out = append(out, svc)
	}
	return out
}

func filterAutoDeployServices(services []model.Service, pushBranch string) []model.Service {
	pushBranch = strings.TrimSpace(pushBranch)
	var out []model.Service
	for _, svc := range services {
		if !svc.AutoDeployEnabled {
			continue
		}
		svcBranch := strings.TrimSpace(svc.Branch)
		if svcBranch == "" {
			svcBranch = "main"
		}
		if svcBranch != pushBranch {
			continue
		}
		out = append(out, svc)
	}
	return out
}
