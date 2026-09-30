package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDefaultWorkspaceDir_PrefersRepoDataWorkspaces(t *testing.T) {
	orig, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(orig) })

	// This test file lives in apps/agent/internal/config → chdir to apps/agent.
	agentDir := filepath.Join(orig, "..", "..")
	agentDir, err = filepath.Abs(agentDir)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(agentDir); err != nil {
		t.Fatal(err)
	}

	got := defaultWorkspaceDir()
	wantSuffix := filepath.Join("data", "workspaces")
	if !strings.HasSuffix(got, wantSuffix) {
		t.Fatalf("want path ending in %q, got %q", wantSuffix, got)
	}
	repoRoot := filepath.Dir(filepath.Dir(got))
	if _, err := os.Stat(filepath.Join(repoRoot, "templates")); err != nil {
		t.Fatalf("expected monorepo templates next to data/workspaces: %v", err)
	}
}
