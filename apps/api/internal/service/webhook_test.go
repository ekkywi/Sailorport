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

type fakeWebhookDeployer struct {
	lastServiceID string
	lastEnv       string
	deployment    model.Deployment
	err           error
	called        bool
}

func (f *fakeWebhookDeployer) Create(ctx context.Context, serviceID string, req model.CreateDeploymentRequest) (model.Deployment, error) {
	f.called = true
	f.lastServiceID = serviceID
	f.lastEnv = req.Environment
	if f.err != nil {
		return model.Deployment{}, f.err
	}
	out := f.deployment
	if out.ID == "" {
		out.ID = "dep-1"
	}
	out.ServiceID = serviceID
	return out, nil
}

func signBody(secret string, body []byte) string {
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write(body)
	return "sha256=" + hex.EncodeToString(mac.Sum(nil))
}

func TestHandleGitHub_PushCreatesDeployment(t *testing.T) {
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
			ID:                    "svc-1",
			SourceType:            "git",
			RepoURL:               "https://github.com/acme/hello.git",
			Branch:                "main",
			WebhookSecret:         secret,
			AutoDeployEnabled:     true,
			AutoDeployEnvironment: "staging",
		}},
	}
	dep := &fakeWebhookDeployer{
		deployment: model.Deployment{ID: "dep-99"},
	}
	ack, err := NewWebhook(cat, dep).HandleGitHub(context.Background(), "push", signBody(secret, body), body)
	if err != nil {
		t.Fatal(err)
	}
	if ack.Ignored || ack.Branch != "main" || ack.Repo != "acme/hello" {
		t.Fatalf("%+v", ack)
	}
	if ack.ServiceID != "svc-1" || ack.DeploymentID != "dep-99" || ack.Environment != "staging" {
		t.Fatalf("ack deploy fields: %+v", ack)
	}
	if !dep.called || dep.lastServiceID != "svc-1" || dep.lastEnv != "staging" {
		t.Fatalf("deployer: called=%v service=%q env=%q", dep.called, dep.lastServiceID, dep.lastEnv)
	}
}

func TestHandleGitHub_IgnoresWhenAutoDeployOff(t *testing.T) {
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
			SourceType:        "git",
			RepoURL:           "https://github.com/acme/hello.git",
			Branch:            "main",
			WebhookSecret:     secret,
			AutoDeployEnabled: false,
		}},
	}
	dep := &fakeWebhookDeployer{}
	ack, err := NewWebhook(cat, dep).HandleGitHub(
		context.Background(), "push", signBody(secret, body), body,
	)
	if err != nil {
		t.Fatal(err)
	}
	if !ack.Ignored {
		t.Fatalf("expected ignored, got %+v", ack)
	}
	if dep.called {
		t.Fatal("deployer should not be called when auto-deploy off")
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
	_, err := NewWebhook(cat, nil).HandleGitHub(context.Background(), "push", "sha256=deadbeef", body)
	if !errors.Is(err, ErrUnauthorized) {
		t.Fatalf("expected ErrUnauthorized, got %v", err)
	}
}

// Repo yang tidak terdaftar harus terlihat sama dengan signature salah, supaya
// endpoint publik ini tidak bisa dipakai mengintip isi catalog.
func TestHandleGitHub_UnknownRepoIsUnauthorized(t *testing.T) {
	body := []byte(`{
		"ref": "refs/heads/main",
		"after": "abc123",
		"repository": {
			"full_name": "acme/hello",
			"clone_url": "https://github.com/acme/hello.git"
		},
		"pusher": {"name": "alice"}
	}`)
	ack, err := NewWebhook(&fakeWebhookCatalog{}, nil).HandleGitHub(context.Background(), "push", "", body)
	if !errors.Is(err, ErrUnauthorized) {
		t.Fatalf("expected ErrUnauthorized, got %v", err)
	}
	if ack.Repo != "" || ack.Reason != "" {
		t.Fatalf("ack should stay empty before auth: %+v", ack)
	}
}

// Dua service satu repo: secret milik service lain tidak boleh memicu deploy.
func TestHandleGitHub_OtherServiceSecretCannotDeploy(t *testing.T) {
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
		services: []model.Service{
			{
				ID:                "svc-a",
				SourceType:        "git",
				RepoURL:           "https://github.com/acme/hello.git",
				Branch:            "main",
				WebhookSecret:     "secret-a",
				AutoDeployEnabled: false,
			},
			{
				ID:                    "svc-b",
				SourceType:            "git",
				RepoURL:               "https://github.com/acme/hello.git",
				Branch:                "main",
				WebhookSecret:         "secret-b",
				AutoDeployEnabled:     true,
				AutoDeployEnvironment: "prod",
			},
		},
	}
	dep := &fakeWebhookDeployer{}
	ack, err := NewWebhook(cat, dep).HandleGitHub(
		context.Background(), "push", signBody("secret-a", body), body,
	)
	if err != nil {
		t.Fatal(err)
	}
	if dep.called {
		t.Fatalf("secret svc-a must not deploy svc-b (deployed %q)", dep.lastServiceID)
	}
	if !ack.Ignored {
		t.Fatalf("expected ignored, got %+v", ack)
	}
}

func TestHandleGitHub_DeploysServiceOwningTheSecret(t *testing.T) {
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
		services: []model.Service{
			{
				ID:                    "svc-a",
				SourceType:            "git",
				RepoURL:               "https://github.com/acme/hello.git",
				Branch:                "main",
				WebhookSecret:         "secret-a",
				AutoDeployEnabled:     true,
				AutoDeployEnvironment: "dev",
			},
			{
				ID:                    "svc-b",
				SourceType:            "git",
				RepoURL:               "https://github.com/acme/hello.git",
				Branch:                "main",
				WebhookSecret:         "secret-b",
				AutoDeployEnabled:     true,
				AutoDeployEnvironment: "prod",
			},
		},
	}
	dep := &fakeWebhookDeployer{}
	ack, err := NewWebhook(cat, dep).HandleGitHub(
		context.Background(), "push", signBody("secret-b", body), body,
	)
	if err != nil {
		t.Fatal(err)
	}
	if dep.lastServiceID != "svc-b" || dep.lastEnv != "prod" {
		t.Fatalf("deployer: service=%q env=%q", dep.lastServiceID, dep.lastEnv)
	}
	if ack.ServiceID != "svc-b" || ack.Environment != "prod" {
		t.Fatalf("ack: %+v", ack)
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
	_, err := NewWebhook(cat, nil).HandleGitHub(context.Background(), "push", "", body)
	if !errors.Is(err, ErrUnauthorized) {
		t.Fatalf("expected ErrUnauthorized, got %v", err)
	}
}

func TestHandleGitHub_IgnoresPing(t *testing.T) {
	ack, err := NewWebhook(nil, nil).HandleGitHub(context.Background(), "ping", "", []byte(`{}`))
	if err != nil || !ack.Ignored {
		t.Fatalf("%+v %v", ack, err)
	}
}

func TestHandleGitHub_MissingEvent(t *testing.T) {
	_, err := NewWebhook(nil, nil).HandleGitHub(context.Background(), "", "", []byte(`{}`))
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
