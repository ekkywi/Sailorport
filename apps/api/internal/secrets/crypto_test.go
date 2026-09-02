package secrets

import (
	"strings"
	"testing"
)

func testSecretsKey(t *testing.T) []byte {
	t.Helper()
	key := make([]byte, secretsKeyBytes)
	for i := range key {
		key[i] = byte(i)
	}
	return key
}

func TestEncryptDecrypt_RoundTrip(t *testing.T) {
	key := testSecretsKey(t)
	plain := "my-super-secret-password"

	enc, err := encryptSecretValue(key, plain)
	if err != nil {
		t.Fatalf("encrypt: %v", err)
	}
	if !isEncryptedValue(enc) {
		t.Fatalf("expected enc:v1: prefix, got %q", enc)
	}
	if enc == plain {
		t.Fatal("ciphertext must differ from plaintext")
	}

	got, err := decryptSecretValue(key, enc)
	if err != nil {
		t.Fatalf("decrypt: %v", err)
	}
	if got != plain {
		t.Fatalf("round-trip: want %q, got %q", plain, got)
	}
}

func TestEncryptDecrypt_DifferentNonceEachTime(t *testing.T) {
	key := testSecretsKey(t)
	a, err := encryptSecretValue(key, "same")
	if err != nil {
		t.Fatal(err)
	}
	b, err := encryptSecretValue(key, "same")
	if err != nil {
		t.Fatal(err)
	}
	if a == b {
		t.Fatal("two encryptions of same plaintext should differ (random nonce)")
	}
	gotA, err := decryptSecretValue(key, a)
	if err != nil {
		t.Fatalf("decrypt a: %v", err)
	}
	gotB, err := decryptSecretValue(key, b)
	if err != nil {
		t.Fatalf("decrypt b: %v", err)
	}
	if gotA != "same" || gotB != "same" {
		t.Fatal("both must decrypt to same plaintext")
	}
}

func TestDecrypt_WrongKey(t *testing.T) {
	key := testSecretsKey(t)
	enc, err := encryptSecretValue(key, "secret")
	if err != nil {
		t.Fatal(err)
	}

	wrongKey := make([]byte, secretsKeyBytes)
	wrongKey[0] = 99
	_, err = decryptSecretValue(wrongKey, enc)
	if err == nil {
		t.Fatal("expected error with wrong key")
	}
}

func TestDecrypt_NotEncrypted(t *testing.T) {
	key := testSecretsKey(t)
	_, err := decryptSecretValue(key, "plaintext-password")
	if err == nil {
		t.Fatal("expected error for plaintext value")
	}
	if !strings.Contains(err.Error(), "not encrypted") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateSecretsKey(t *testing.T) {
	if err := validateSecretsKey(make([]byte, 16)); err == nil {
		t.Fatal("expected error for 16-byte key")
	}
	if err := validateSecretsKey(make([]byte, secretsKeyBytes)); err != nil {
		t.Fatalf("valid key: %v", err)
	}
}
