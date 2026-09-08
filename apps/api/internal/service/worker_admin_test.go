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
		Tier:         "nonprod",
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

func TestValidateWorkerTier(t *testing.T) {
	if err := validateWorkerTier(""); err != nil {
		t.Fatalf("empty tier OK: %v", err)
	}
	if err := validateWorkerTier("nonprod"); err != nil {
		t.Fatalf("nonprod OK: %v", err)
	}
	if err := validateWorkerTier("prod"); err != nil {
		t.Fatalf("prod OK: %v", err)
	}
	if err := validateWorkerTier("gold"); err == nil {
		t.Fatal("expected error for unknown tier")
	}
}

func TestNormalizeWorkerEnvironmentsInput(t *testing.T) {
	got := normalizeWorkerEnvironmentsInput(" Dev, staging,dev , ")
	if len(got) != 2 || got[0] != "dev" || got[1] != "staging" {
		t.Fatalf("got %#v", got)
	}
	if len(normalizeWorkerEnvironmentsInput("")) != 0 {
		t.Fatal("empty should yield no slugs")
	}
}
