package service

import (
	"errors"
	"testing"

	"github.com/ekkywi/sailorport/apps/api/internal/model"
)

func TestWorkerAllowsEnvironment(t *testing.T) {
	cases := []struct {
		name   string
		labels map[string]any
		env    string
		want   bool
	}{
		{"unlabeled allows prod", nil, "prod", true},
		{"empty environments allows staging", map[string]any{"role": "agent"}, "staging", true},
		{"nonprod allows dev", map[string]any{"environments": "dev,staging"}, "dev", true},
		{"nonprod allows staging", map[string]any{"environments": "dev,staging"}, "staging", true},
		{"nonprod rejects prod", map[string]any{"environments": "dev,staging"}, "prod", false},
		{"spaces and case", map[string]any{"environments": " Dev , Staging "}, "staging", true},
		{"prod only", map[string]any{"environments": "prod"}, "dev", false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			w := model.Worker{Name: "w1", Labels: tc.labels}
			got := workerAllowsEnvironment(w, tc.env)
			if got != tc.want {
				t.Fatalf("got %v want %v", got, tc.want)
			}
		})
	}
}

func TestWorkerEnvConflict_WrapsErrConflict(t *testing.T) {
	w := model.Worker{
		Name:   "nonprod-01",
		Labels: map[string]any{"environments": "dev, staging"},
	}
	err := workerEnvConflict(w, "prod")
	if !errors.Is(err, ErrConflict) {
		t.Fatalf("expected ErrConflict, got %v", err)
	}
}
