package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ekkywi/sailorport/apps/api/internal/model"
)

// webhookCatalog is the catalog surface needed for webhook matching.
type webhookCatalog interface {
	List(ctx context.Context) ([]model.Service, error)
}

type webhookDeployer interface {
	Create(ctx context.Context, serviceID string, req model.CreateDeploymentRequest) (model.Deployment, error)
}

type Webhook struct {
	catalog     webhookCatalog
	deployments webhookDeployer
}

func NewWebhook(catalog webhookCatalog, deployments webhookDeployer) *Webhook {
	return &Webhook{
		catalog:     catalog,
		deployments: deployments,
	}
}

// HandleGitHub verifies a GitHub push webhook and may create a deployment when auto-deploy is on.
//
// Signature diverifikasi sebelum apa pun yang bergantung isi catalog ikut terjawab,
// dan semua kegagalan auth (repo tidak dikenal, service tanpa secret, signature
// salah) memakai satu error yang sama — endpoint ini publik, jadi bedanya balasan
// bisa dipakai orang luar untuk mengintip repo mana yang terdaftar.
func (w *Webhook) HandleGitHub(
	ctx context.Context,
	event string,
	signatureHeader string,
	body []byte,
) (model.WebhookAck, error) {
	event = strings.TrimSpace(event)
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

	dep, err := w.deployments.Create(ctx, target.ID, model.CreateDeploymentRequest{
		Environment: env,
	})
	if err != nil {
		return model.WebhookAck{}, err
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
	all, err := w.catalog.List(ctx)
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
