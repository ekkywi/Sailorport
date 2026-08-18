package service

import (
	"testing"
	"time"

	"github.com/ekkywi/sailorport/apps/api/internal/model"
)

func TestLatestPerEnv_PicksNewestPerSlug(t *testing.T) {
	now := time.Now()
	deps := []model.Deployment{
		{EnvironmentSlug: "dev", Status: "stopped", CreatedAt: now},
		{EnvironmentSlug: "staging", Status: "running", CreatedAt: now},
		{EnvironmentSlug: "dev", Status: "running", CreatedAt: now.Add(-time.Hour)},
	}
	out := latestPerEnv(deps)
	if out["dev"].Status != "stopped" {
		t.Fatalf("expected newest dev, got %s", out["dev"].Status)
	}
	if out["staging"].Status != "running" {
		t.Fatalf("expected staging running")
	}
	if len(out) != 2 {
		t.Fatalf("expected 2 envs, got %d", len(out))
	}
}

func TestLatestPerEnv_EmptySlugDefaultsToDev(t *testing.T) {
	out := latestPerEnv([]model.Deployment{
		{EnvironmentSlug: "  ", Status: "failed"},
	})
	if got, ok := out["dev"]; !ok || got.Status != "failed" {
		t.Fatalf("expected empty slug mapped to dev, got %+v", out)
	}
}
