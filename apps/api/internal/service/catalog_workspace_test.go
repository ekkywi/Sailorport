package service

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/ekkywi/sailorport/apps/api/internal/model"
)

func TestSanitizeCreateWorkspacePath_ClearsNonScaffold(t *testing.T) {
	c := &Catalog{workspaceDir: t.TempDir()}
	req := model.CreateServiceRequest{
		SourceType:    "git",
		WorkspacePath: "/etc/passwd",
	}
	if err := c.sanitizeCreateWorkspacePath(&req); err != nil {
		t.Fatal(err)
	}
	if req.WorkspacePath != "" {
		t.Fatalf("want empty, got %q", req.WorkspacePath)
	}
}

func TestSanitizeCreateWorkspacePath_AllowsUnderRoot(t *testing.T) {
	root := t.TempDir()
	c := &Catalog{workspaceDir: root}
	child := filepath.Join(root, "my-svc")
	req := model.CreateServiceRequest{
		SourceType:    "scaffold",
		WorkspacePath: child,
	}
	if err := c.sanitizeCreateWorkspacePath(&req); err != nil {
		t.Fatal(err)
	}
	if req.WorkspacePath != child {
		t.Fatalf("got %q", req.WorkspacePath)
	}
}

func TestSanitizeCreateWorkspacePath_RejectsEscape(t *testing.T) {
	root := t.TempDir()
	c := &Catalog{workspaceDir: root}
	req := model.CreateServiceRequest{
		SourceType:    "scaffold",
		WorkspacePath: filepath.Join(filepath.Dir(root), "outside"),
	}
	err := c.sanitizeCreateWorkspacePath(&req)
	if !errors.Is(err, ErrInvalid) {
		t.Fatalf("want ErrInvalid, got %v", err)
	}
}
