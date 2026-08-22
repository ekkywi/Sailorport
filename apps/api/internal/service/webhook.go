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

type Webhook struct {
	catalog webhookCatalog
}

func NewWebhook(catalog webhookCatalog) *Webhook {
	return &Webhook{catalog: catalog}
}

// HandleGitHub acknowledges a GitHub webhook after matching a service and verifying HMAC.
// Creating a deployment is Step 20d.
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

	ack.Repo = payload.Repository.FullName
	ack.CloneURL = payload.Repository.CloneURL
	ack.Branch = branch
	ack.CommitSHA = payload.After
	ack.Pusher = payload.Pusher.Name

	matches, err := w.findServicesByCloneURL(ctx, payload.Repository.CloneURL)
	if err != nil {
		return model.WebhookAck{}, err
	}
	if len(matches) == 0 {
		ack.Ignored = true
		ack.Reason = "no matching service"
		return ack, nil
	}

	secret := ""
	for _, svc := range matches {
		if s := strings.TrimSpace(svc.WebhookSecret); s != "" {
			secret = s
			break
		}
	}
	if secret == "" {
		return model.WebhookAck{}, fmt.Errorf("%w: webhook secret not configured", ErrUnauthorized)
	}

	if err := verifyGitHubSignature(secret, body, signatureHeader); err != nil {
		return model.WebhookAck{}, err
	}

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
