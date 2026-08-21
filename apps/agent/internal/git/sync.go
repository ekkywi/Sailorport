package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func Sync(repoURL, branch, dir string) error {
	repoURL = strings.TrimSpace(repoURL)
	branch = strings.TrimSpace(branch)
	dir = strings.TrimSpace(dir)

	if repoURL == "" {
		return fmt.Errorf("repo URL is empty")
	}
	if branch == "" {
		branch = "main"
	}
	if dir == "" {
		return fmt.Errorf("target dir is empty")
	}

	gitDir := filepath.Join(dir, ".git")
	if _, err := os.Stat(gitDir); err == nil {
		return pull(dir, branch)
	}
	if err := os.MkdirAll(filepath.Dir(dir), 0o755); err != nil {
		return fmt.Errorf("mkdir parent: %w", err)
	}
	_ = os.RemoveAll(dir)

	cmd := exec.Command("git", "clone", "--branch", branch, "--single-branch", repoURL, dir)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git clone: %w\n%s", err, out)
	}
	return nil
}

func pull(dir, branch string) error {
	cmds := [][]string{
		{"fetch", "origin"},
		{"checkout", branch},
		{"pull", "origin", branch},
	}
	for _, args := range cmds {
		cmd := exec.Command("git", args...)
		cmd.Dir = dir
		out, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("git %v: %w\n%s", args, err, out)
		}
	}
	return nil
}
