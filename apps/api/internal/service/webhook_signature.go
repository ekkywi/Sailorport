package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

func verifyGitHubSignature(secret string, body []byte, signatureHeader string) error {
	secret = strings.TrimSpace(secret)
	if secret == "" {
		return fmt.Errorf("%w: webhook secret not configured", ErrUnauthorized)
	}

	sig := strings.TrimSpace(signatureHeader)
	const prefix = "sha256="
	if !strings.HasPrefix(sig, prefix) {
		return fmt.Errorf("%w: missing or invalid signature header", ErrUnauthorized)
	}
	wantHex := strings.TrimPrefix(sig, prefix)

	want, err := hex.DecodeString(wantHex)
	if err != nil {
		return fmt.Errorf("%w: invalid signature encoding", ErrUnauthorized)
	}

	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write(body)
	got := mac.Sum(nil)

	if !hmac.Equal(got, want) {
		return fmt.Errorf("%w: invalid signature", ErrUnauthorized)
	}
	return nil
}
