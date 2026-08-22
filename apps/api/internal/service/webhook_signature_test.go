package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"testing"
)

func TestVerifyGitHubSignature_OK(t *testing.T) {
	secret := "my-secret"
	body := []byte(`{"ref":"refs/heads/main"}`)
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write(body)
	header := "sha256=" + hex.EncodeToString(mac.Sum(nil))

	if err := verifyGitHubSignature(secret, body, header); err != nil {
		t.Fatal(err)
	}
}

func TestVerifyGitHubSignature_Bad(t *testing.T) {
	err := verifyGitHubSignature("my-secret", []byte(`{}`), "sha256=00")
	if !errors.Is(err, ErrUnauthorized) {
		t.Fatalf("got %v", err)
	}
}

func TestVerifyGitHubSignature_EmptySecret(t *testing.T) {
	err := verifyGitHubSignature("", []byte(`{}`), "sha256=ab")
	if !errors.Is(err, ErrUnauthorized) {
		t.Fatalf("got %v", err)
	}
}
