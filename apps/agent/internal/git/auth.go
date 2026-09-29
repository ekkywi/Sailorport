package git

import (
	"fmt"
	"net/url"
	"os/exec"
	"strings"
)

func resolveAuthURL(repoURL, token string) (authURL, cleanURL string, err error) {
	repoURL = strings.TrimSpace(repoURL)
	token = strings.TrimSpace(token)
	cleanURL = repoURL

	if token == "" {
		return repoURL, cleanURL, nil
	}

	lower := strings.ToLower(repoURL)
	if !strings.HasPrefix(lower, "https://") && !strings.HasPrefix(lower, "http://") {
		return "", "", fmt.Errorf("git_token requires http(s) repo_url (not ssh)")
	}

	u, err := url.Parse(repoURL)
	if err != nil {
		return "", "", fmt.Errorf("parse repo url: %w", err)
	}
	if u.Host == "" {
		return "", "", fmt.Errorf("repo url missing host")
	}

	u.User = url.UserPassword("x-access-token", token)
	return u.String(), cleanURL, nil
}

func setRemoteURL(dir, remoteURL string) error {
	cmd := exec.Command("git", "remote", "set-url", "origin", "--", remoteURL)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git remote set-url: %w\n%s", err, out)
	}
	return nil
}
