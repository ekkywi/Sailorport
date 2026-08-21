package service

import (
	"errors"
	"strings"
	"testing"

	"github.com/ekkywi/sailorport/apps/api/internal/model"
)

func TestNormalizeCreate_RequiresName(t *testing.T) {
	_, err := normalizeCreate(model.CreateServiceRequest{
		Name:        "   ",
		Description: "x",
		Owner:       "y",
	})
	if !errors.Is(err, ErrInvalid) {
		t.Fatalf("expected ErrInvalid, got %v", err)
	}
	if !strings.Contains(err.Error(), "name is required") {
		t.Fatalf("unexpected message: %v", err)
	}
}

func TestNormalizeCreate_TrimsFields(t *testing.T) {
	req, err := normalizeCreate(model.CreateServiceRequest{
		Name:        "  payments-api  ",
		Description: "  desc  ",
		Owner:       "  team  ",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if req.Name != "payments-api" || req.Description != "desc" || req.Owner != "team" {
		t.Fatalf("fields not trimmed: %+v", req)
	}
	if req.SourceType != "scaffold" || req.Branch != "main" || req.DockerfilePath != "Dockerfile" {
		t.Fatalf("defaults not applied: source=%q branch=%q dockerfile=%q",
			req.SourceType, req.Branch, req.DockerfilePath)
	}
}

func TestNormalizeCreate_GitRequiresRepoURL(t *testing.T) {
	_, err := normalizeCreate(model.CreateServiceRequest{
		Name:       "from-git",
		SourceType: "git",
	})
	if !errors.Is(err, ErrInvalid) {
		t.Fatalf("expected ErrInvalid, got %v", err)
	}
	if !strings.Contains(err.Error(), "repo_url") {
		t.Fatalf("unexpected message: %v", err)
	}
}

func TestNormalizeCreate_GitOK(t *testing.T) {
	req, err := normalizeCreate(model.CreateServiceRequest{
		Name:       "from-git",
		SourceType: "git",
		RepoURL:    "https://github.com/example/hello.git",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if req.SourceType != "git" || req.RepoURL == "" || req.Branch != "main" {
		t.Fatalf("unexpected: %+v", req)
	}
}

func TestNormalizeCreate_InvalidSourceType(t *testing.T) {
	_, err := normalizeCreate(model.CreateServiceRequest{
		Name:       "x",
		SourceType: "ftp",
	})
	if !errors.Is(err, ErrInvalid) {
		t.Fatalf("expected ErrInvalid, got %v", err)
	}
}

func TestNormalizeCreate_WebhookDefaults(t *testing.T) {
	req, err := normalizeCreate(model.CreateServiceRequest{
		Name: "payments-api",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if req.AutoDeployEnvironment != "staging" {
		t.Fatalf("want staging, got %q", req.AutoDeployEnvironment)
	}
	if req.AutoDeployEnabled {
		t.Fatal("enabled should default false")
	}
	if req.WebhookSecret != "" {
		t.Fatalf("secret should default empty, got %q", req.WebhookSecret)
	}
}

func TestNormalizeCreate_InvalidAutoDeployEnv(t *testing.T) {
	_, err := normalizeCreate(model.CreateServiceRequest{
		Name:                  "x",
		AutoDeployEnvironment: "qa",
	})
	if !errors.Is(err, ErrInvalid) {
		t.Fatalf("expected ErrInvalid, got %v", err)
	}
}

func TestNormalizeUpdate_MergesGitFields(t *testing.T) {
	existing := model.Service{
		SourceType:     "git",
		RepoURL:        "https://github.com/example/hello.git",
		Branch:         "develop",
		DockerfilePath: "deploy/Dockerfile",
	}
	req, err := normalizeUpdate(model.UpdateServiceRequest{
		Name:        "from-git",
		Description: "updated",
		Owner:       "team",
	}, existing)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if req.SourceType != "git" || req.RepoURL != existing.RepoURL {
		t.Fatalf("merge failed: %+v", req)
	}
	if req.Branch != "develop" || req.DockerfilePath != "deploy/Dockerfile" {
		t.Fatalf("branch/dockerfile merge failed: %+v", req)
	}
}

func TestNormalizeUpdate_SwitchToGitRequiresRepo(t *testing.T) {
	existing := model.Service{SourceType: "scaffold"}
	_, err := normalizeUpdate(model.UpdateServiceRequest{
		Name:       "x",
		SourceType: "git",
	}, existing)
	if !errors.Is(err, ErrInvalid) {
		t.Fatalf("expected ErrInvalid, got %v", err)
	}
}

func TestNormalizeUpdate_MergesWebhookFields(t *testing.T) {
	existing := model.Service{
		SourceType:            "scaffold",
		WebhookSecret:         "super-secret",
		AutoDeployEnabled:     true,
		AutoDeployEnvironment: "prod",
	}
	req, err := normalizeUpdate(model.UpdateServiceRequest{
		Name:        "payments-api",
		Description: "updated",
		Owner:       "team",
	}, existing)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if req.WebhookSecret != "super-secret" {
		t.Fatalf("secret wiped: %q", req.WebhookSecret)
	}
	if !req.AutoDeployEnabled {
		t.Fatal("enabled wiped")
	}
	if req.AutoDeployEnvironment != "prod" {
		t.Fatalf("env wiped: %q", req.AutoDeployEnvironment)
	}
}

func TestNormalizeUpdate_InvalidAutoDeployEnv(t *testing.T) {
	existing := model.Service{
		SourceType: "scaffold",
		AutoDeployEnvironment: "staging",
	}
	_, err := normalizeUpdate(model.UpdateServiceRequest{
		Name: "x",
		AutoDeployEnvironment: "qa",
	}, existing)
	if !errors.Is(err, ErrInvalid) {
		t.Fatalf("expected ErrInvalid, got %v", err)
	}
}