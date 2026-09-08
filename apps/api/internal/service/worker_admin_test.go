package service

import (
	"testing"

	"github.com/ekkywi/sailorport/apps/api/internal/model"
)

func TestMergeWorkerCapabilityLabels(t *testing.T) {
	existing := map[string]any{
		"role": "agent",
		"tier": "old",
	}
	got := mergeWorkerCapabilityLabels(existing, model.UpdateWorkerLabelsRequest{
		Tier: "nonprod",
		Environments: "dev,staging",
	})
	if got["role"] != "agent" {
		t.Fatalf("role should be preserved, got %#v", got["role"])
	}
	if got["tier"] != "nonprod" {
		t.Fatalf("tier: got %#v", got["tier"])
	}
	if got["environments"] != "dev,staging" {
		t.Fatalf("environments: got %#v", got["environments"])
	}
	if existing["tier"] != "old" {
		t.Fatalf("merge must not mutate existing map")
	}
}