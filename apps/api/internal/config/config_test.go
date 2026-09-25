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
	hexKey := strings.Repeat("ab", 32)
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
	cfg := Config{SecretsKey: "abcd"}
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

func TestParseCORSOrigins_Empty(t *testing.T) {
	if parseCORSOrigins("") != nil {
		t.Fatal("expected nil")
	}
	if parseCORSOrigins(" , , ") != nil {
		t.Fatal("expected nil for blank parts")
	}
}

func TestParseCORSOrigins_List(t *testing.T) {
	got := parseCORSOrigins(" http://localhost:5173, http://127.0.0.1:5173 ")
	if len(got) != 2 {
		t.Fatalf("got %#v", got)
	}
	if got[0] != "http://localhost:5173" || got[1] != "http://127.0.0.1:5173" {
		t.Fatalf("got %#v", got)
	}
}

func TestLoad_CORSOrigins_DevDefault(t *testing.T) {
	t.Setenv("APP_ENV", "development")
	t.Setenv("CORS_ORIGINS", "")
	cfg := Load()
	if len(cfg.CORSOrigins) != 2 {
		t.Fatalf("dev default: %#v", cfg.CORSOrigins)
	}
}

func TestLoad_CORSOrigins_ProductionEmpty(t *testing.T) {
	t.Setenv("APP_ENV", "production")
	t.Setenv("CORS_ORIGINS", "")
	cfg := Load()
	if cfg.CORSOrigins != nil && len(cfg.CORSOrigins) != 0 {
		t.Fatalf("production empty should be no origins, got %#v", cfg.CORSOrigins)
	}
}

func TestLoad_CORSOrigins_FromEnv(t *testing.T) {
	t.Setenv("APP_ENV", "production")
	t.Setenv("CORS_ORIGINS", "https://portal.example.com")
	cfg := Load()
	if len(cfg.CORSOrigins) != 1 || cfg.CORSOrigins[0] != "https://portal.example.com" {
		t.Fatalf("got %#v", cfg.CORSOrigins)
	}
}