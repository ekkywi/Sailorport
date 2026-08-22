package service

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"testing"

	"github.com/ekkywi/sailorport/apps/api/internal/model"
)

type fakeWebhookCatalog struct {
	services []model.Service
	err      error
}

func (f *fakeWebhookCatalog) List(ctx context.Context) ([]model.Service, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.services, nil
}

func signBody(secret string, body []byte) string {
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write(body)
	return "sha256=" + hex.EncodeToString(mac.Sum(nil))
}

func TestHandleGitHub_PushSigned(t *testing.T) {
	secret := "test-secret"
	body := []byte(`{
		"ref": "refs/heads/main",
		"after": "abc123",
		"repository": {
			"full_name": "acme/hello",
			"clone_url": "https://github.com/acme/hello.git"
		},
		"pusher": {"name": "alice"}
	}`)
	cat := &fakeWebhookCatalog{
		services: []model.Service{{
			SourceType:    "git",
			RepoURL:       "https://github.com/acme/hello.git",
			WebhookSecret: secret,
		}},
	}
	ack, err := NewWebhook(cat).HandleGitHub(context.Background(), "push", signBody(secret, body), body)
	if err != nil {
		t.Fatal(err)
	}
	if ack.Ignored || ack.Branch != "main" || ack.Repo != "acme/hello" {
		t.Fatalf("%+v", ack)
	}
	if ack.CloneURL != "https://github.com/acme/hello.git" || ack.CommitSHA != "abc123" {
		t.Fatalf("%+v", ack)
	}
}

func TestHandleGitHub_BadSignature(t *testing.T) {
	body := []byte(`{
		"ref": "refs/heads/main",
		"after": "abc123",
		"repository": {
			"full_name": "acme/hello",
			"clone_url": "https://github.com/acme/hello.git"
		},
		"pusher": {"name": "alice"}
	}`)
	cat := &fakeWebhookCatalog{
		services: []model.Service{{
			SourceType:    "git",
			RepoURL:       "https://github.com/acme/hello",
			WebhookSecret: "test-secret",
		}},
	}
	_, err := NewWebhook(cat).HandleGitHub(context.Background(), "push", "sha256=deadbeef", body)
	if !errors.Is(err, ErrUnauthorized) {
		t.Fatalf("expected ErrUnauthorized, got %v", err)
	}
}

func TestHandleGitHub_NoMatchingService(t *testing.T) {
	body := []byte(`{
		"ref": "refs/heads/main",
		"after": "abc123",
		"repository": {
			"full_name": "acme/hello",
			"clone_url": "https://github.com/acme/hello.git"
		},
		"pusher": {"name": "alice"}
	}`)
	ack, err := NewWebhook(&fakeWebhookCatalog{}).HandleGitHub(context.Background(), "push", "", body)
	if err != nil {
		t.Fatal(err)
	}
	if !ack.Ignored || ack.Reason != "no matching service" {
		t.Fatalf("%+v", ack)
	}
}

func TestHandleGitHub_SecretNotConfigured(t *testing.T) {
	body := []byte(`{
		"ref": "refs/heads/main",
		"after": "abc123",
		"repository": {
			"full_name": "acme/hello",
			"clone_url": "https://github.com/acme/hello.git"
		},
		"pusher": {"name": "alice"}
	}`)
	cat := &fakeWebhookCatalog{
		services: []model.Service{{
			SourceType: "git",
			RepoURL:    "https://github.com/acme/hello.git",
		}},
	}
	_, err := NewWebhook(cat).HandleGitHub(context.Background(), "push", "", body)
	if !errors.Is(err, ErrUnauthorized) {
		t.Fatalf("expected ErrUnauthorized, got %v", err)
	}
}

func TestHandleGitHub_IgnoresPing(t *testing.T) {
	ack, err := NewWebhook(nil).HandleGitHub(context.Background(), "ping", "", []byte(`{}`))
	if err != nil || !ack.Ignored {
		t.Fatalf("%+v %v", ack, err)
	}
}

func TestHandleGitHub_MissingEvent(t *testing.T) {
	_, err := NewWebhook(nil).HandleGitHub(context.Background(), "", "", []byte(`{}`))
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestBranchFromRef(t *testing.T) {
	if branchFromRef("refs/heads/develop") != "develop" {
		t.Fatal("branch")
	}
	if branchFromRef("refs/tags/v1") != "" {
		t.Fatal("tag should be empty")
	}
}

func TestNormalizeRepoURL(t *testing.T) {
	a := normalizeRepoURL("https://github.com/Acme/Hello.git")
	b := normalizeRepoURL("https://github.com/acme/hello")
	if a != b {
		t.Fatalf("%q vs %q", a, b)
	}
}
