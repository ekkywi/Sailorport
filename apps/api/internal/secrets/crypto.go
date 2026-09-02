package secrets

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"strings"
)

const encryptedValuePrefix = "enc:1"

const secretsKeyBytes = 32

func isEncryptedValue(stored string) bool {
	return strings.HasPrefix(stored, encryptedValuePrefix)
}

func validateSecretsKey(key []byte) error {
	if len(key) != secretsKeyBytes {
		return fmt.Errorf("secrets key must be %d bytes, got %d", secretsKeyBytes, len(key))
	}
	return nil
}

func encryptSecretValue(key []byte, plaintext string) (string, error) {
	if err := validateSecretsKey(key); err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("aes cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("gcm: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("nonce: %w", err)
	}

	sealed := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	encoded := base64.StdEncoding.EncodeToString(sealed)
	return encryptedValuePrefix + encoded, nil
}

func decryptSecretValue(key []byte, stored string) (string, error) {
	if err := validateSecretsKey(key); err != nil {
		return "", err
	}
	if !isEncryptedValue(stored) {
		return "", fmt.Errorf("value is not encrypted")
	}

	raw, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(stored, encryptedValuePrefix))
	if err != nil {
		return "", fmt.Errorf("base64 decode: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("aes cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("gcm: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(raw) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := raw[:nonceSize], raw[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("decrypt: %w", err)
	}

	return string(plaintext), nil
}
