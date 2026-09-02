package config

import (
	"strings"
	"testing"
)

func TestSecretsKeyBytes_Empty(t *testing.T) {
	key, err := Config{}.SecretsKeyBytes()
	if err != nil {
		t.Fatal(err)
	}
	if key != nil {
		t.Fatalf("expected nil key, got %v", key)
	}
}

func TestSecretsKeyBytes_ValidHex(t *testing.T) {
	hexKey := strings.Repeat("ab", 32) // 64 hex chars = 32 bytes
	cfg := Config{SecretsKey: hexKey}
	key, err := cfg.SecretsKeyBytes()
	if err != nil {
		t.Fatal(err)
	}
	if len(key) != 32 {
		t.Fatalf("len %d", len(key))
	}
}

func TestSecretsKeyBytes_InvalidHex(t *testing.T) {
	cfg := Config{SecretsKey: "not-hex"}
	_, err := cfg.SecretsKeyBytes()
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestSecretsKeyBytes_WrongLength(t *testing.T) {
	cfg := Config{SecretsKey: "abcd"} // 2 bytes
	_, err := cfg.SecretsKeyBytes()
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "32 bytes") {
		t.Fatalf("unexpected: %v", err)
	}
}

func TestValidate_Development_AllowsMissingSecretsKey(t *testing.T) {
	cfg := Config{
		AppEnv:     "development",
		JWTSecret:  devJWTSecret,
		AgentToken: devAgentToken,
		SecretsKey: "",
	}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("development should allow missing secrets key: %v", err)
	}
}

func TestValidate_Production_RequiresSecretsKey(t *testing.T) {
	cfg := Config{
		AppEnv:     "production",
		JWTSecret:  "real-jwt-secret-at-least-32-chars-long",
		AgentToken: "real-agent-token",
		SecretsKey: "",
	}
	err := cfg.Validate()
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "SAILORPORT_SECRETS_KEY") {
		t.Fatalf("unexpected: %v", err)
	}
}
