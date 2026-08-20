package config

import "testing"

func TestParseWorkerLabels_DefaultRole(t *testing.T) {
	t.Setenv("SAILORPORT_WORKER_LABELS", "")
	t.Setenv("SAILORPORT_WORKER_TIER", "")
	t.Setenv("SAILORPORT_WORKER_ENVIRONMENTS", "")

	labels := parseWorkerLabels()
	if labels["role"] != "agent" {
		t.Fatalf("expected role=agent, got %v", labels["role"])
	}
}

func TestParseWorkerLabels_ConvenienceEnv(t *testing.T) {
	t.Setenv("SAILORPORT_WORKER_TIER", "nonprod")
	t.Setenv("SAILORPORT_WORKER_ENVIRONMENTS", "dev,staging")

	labels := parseWorkerLabels()
	if labels["tier"] != "nonprod" {
		t.Fatalf("tier: got %v", labels["tier"])
	}
	if labels["environments"] != "dev,staging" {
		t.Fatalf("environments: got %v", labels["environments"])
	}
}
