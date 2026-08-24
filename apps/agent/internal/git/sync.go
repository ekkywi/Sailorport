package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func Sync(repoURL, branch, dir, sha string) error {
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
		if err := pull(dir, branch); err != nil {
			return err
		}
	} else {
		if err := os.MkdirAll(filepath.Dir(dir), 0o755); err != nil {
			return fmt.Errorf("mkdir parent: %w", err)
		}
		_ = os.RemoveAll(dir)

		cmd := exec.Command("git", "clone", "--branch", branch, "--single-branch", repoURL, dir)
		out, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("git clone: %w\n%s", err, out)
		}
	}

	sha = strings.TrimSpace(sha)
	if sha == "" {
		return nil
	}
	return checkoutSHA(dir, sha)
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

func HeadSHA(dir string) (string, error) {
	dir = strings.TrimSpace(dir)
	if dir == "" {
		return "", fmt.Errorf("target dir is empty")
	}
	cmd := exec.Command("git", "rev-parse", "HEAD")
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("git rev-parse HEAD: %w\n%s", err, out)
	}
	sha := strings.TrimSpace(string(out))
	if sha == "" {
		return "", fmt.Errorf("empty git sha")
	}
	return sha, nil
}

func checkoutSHA(dir, sha string) error {
	sha = strings.TrimSpace(sha)
	if sha == "" {
		return fmt.Errorf("git sha is empty")
	}

	// Prefer fetching the exact object (helps with --single-branch clones).
	fetch := exec.Command("git", "fetch", "origin", sha)
	fetch.Dir = dir
	if out, err := fetch.CombinedOutput(); err != nil {
		fetchAll := exec.Command("git", "fetch", "origin")
		fetchAll.Dir = dir
		if out2, err2 := fetchAll.CombinedOutput(); err2 != nil {
			return fmt.Errorf("git fetch sha: %w\n%s\n%s", err, out, out2)
		}
	}

	cmd := exec.Command("git", "checkout", "--detach", sha)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git checkout %s: %w\n%s", sha, err, out)
	}
	return nil
}