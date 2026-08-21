package service

import "testing"

func TestHandleGitHub_Push(t *testing.T) {
	body := []byte(`{
		"ref": "refs/heads/main",
		"after": "abc123",
		"repository": {
			"full_name": "acme/hello",
			"clone_url": "https://github.com/acme/hello.git"
		},
		"pusher": {"name": "alice"}
	}`)
	ack, err := NewWebhook().HandleGitHub("push", body)
	if err != nil {
		t.Fatal(err)
	}
	if ack.Ignored || ack.Branch != "main" || ack.Repo != "acme/hello" {
		t.Fatalf("%+v", ack)
	}
	if ack.CloneURL != "https://github.com/acme/hello.git" || ack.CommitSHA != "abc123" {
		t.Fatalf("%+v", ack)
	}
}

func TestHandleGitHub_IgnoresPing(t *testing.T) {
	ack, err := NewWebhook().HandleGitHub("ping", []byte(`{}`))
	if err != nil || !ack.Ignored {
		t.Fatalf("%+v %v", ack, err)
	}
}

func TestHandleGitHub_MissingEvent(t *testing.T) {
	_, err := NewWebhook().HandleGitHub("", []byte(`{}`))
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestBranchFromRef(t *testing.T) {
	if branchFromRef("refs/heads/develop") != "develop" {
		t.Fatal("branch")
	}
	if branchFromRef("refs/tags/v1") != "" {
		t.Fatal("tag should be empty")
	}
}
