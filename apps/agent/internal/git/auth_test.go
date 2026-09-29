package git

import (
	"strings"
	"testing"
)

func TestResolveAuthURL_EmptyToken(t *testing.T) {
	auth, clean, err := resolveAuthURL("https://github.com/acme/r.git", "")
	if err != nil {
		t.Fatal(err)
	}
	if auth != clean || auth != "https://github.com/acme/r.git" {
		t.Fatalf("auth=%q clean=%q", auth, clean)
	}
}

func TestResolveAuthURL_InjectsToken(t *testing.T) {
	auth, clean, err := resolveAuthURL("https://github.com/acme/r.git", "ghp_secret")
	if err != nil {
		t.Fatal(err)
	}
	if clean != "https://github.com/acme/r.git" {
		t.Fatalf("clean: %q", clean)
	}
	if !strings.Contains(auth, "ghp_secret") || !strings.Contains(auth, "x-access-token") {
		t.Fatalf("auth missing creds: %q", auth)
	}
	if strings.Contains(clean, "ghp_secret") {
		t.Fatal("token leaked into clean URL")
	}
}

func TestResolveAuthURL_RejectsSSHWithToken(t *testing.T) {
	_, _, err := resolveAuthURL("git@github.com:acme/r.git", "ghp_x")
	if err == nil {
		t.Fatal("expected error")
	}
}
