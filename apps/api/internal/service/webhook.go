package service

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ekkywi/sailorport/apps/api/internal/model"
)

type Webhook struct{}

func NewWebhook() *Webhook {
	return &Webhook{}
}

// HandleGitHub acknowledges a GitHub webhook. Signature check and deploy follow in 20c/20d.
func (w *Webhook) HandleGitHub(event string, body []byte) (model.WebhookAck, error) {
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
	return ack, nil
}

func branchFromRef(ref string) string {
	const prefix = "refs/heads/"
	if !strings.HasPrefix(ref, prefix) {
		return ""
	}
	return strings.TrimPrefix(ref, prefix)
}
