package service

import (
	"testing"

	"github.com/ekkywi/sailorport/apps/api/internal/model"
)

func TestPublicService_RedactsSecret(t *testing.T) {
	svc := model.Service{
		Name:          "demo",
		WebhookSecret: "super-secret",
	}
	out := PublicService(svc)
	if out.WebhookSecret != "" {
		t.Fatalf("secret leaked: %q", out.WebhookSecret)
	}
	if !out.WebhookSecretSet {
		t.Fatal("expected webhook_secret_set true")
	}
	if svc.WebhookSecret != "super-secret" {
		t.Fatal("input should stay intact (value copy)")
	}
}

func TestPublicService_EmptySecret(t *testing.T) {
	out := PublicService(model.Service{Name: "demo"})
	if out.WebhookSecretSet {
		t.Fatal("expected webhook_secret_set false")
	}
}
