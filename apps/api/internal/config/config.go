package config

import (
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	devJWTSecret  = "dev-only-change-me"
	devAgentToken = "dev-agent-token"
)

type Config struct {
	Port          string
	AppEnv        string
	Version       string
	DatabaseURL   string
	WorkspaceDir  string
	TemplatesDir  string
	JWTSecret     string
	AgentToken    string
	CatalogAppDir string
	SecretsKey    string
}

func Load() Config {
	port := getenv("PORT", "8080")
	appEnv := getenv("APP_ENV", "development")
	version := getenv("APP_VERSION", "0.1.0")
	databaseURL := getenv(
		"DATABASE_URL",
		"postgres://sailorport:sailorport@localhost:5433/sailorport?sslmode=disable",
	)
	workspaceDir := getenv("SAILORPORT_WORKSPACE", defaultWorkspaceDir())
	templatesDir := getenv("SAILORPORT_TEMPLATES", defaultTemplatesDir())
	jwtSecret := getenv("AUTH_JWT_SECRET", devJWTSecret)
	agentToken := getenv("SAILORPORT_AGENT_TOKEN", devAgentToken)
	catalogAppsDir := getenv("SAILORPORT_CATALOG_APPS", defaultCatalogAppsDir())
	secretsKey := getenv("SAILORPORT_SECRETS_KEY", "")

	return Config{
		Port:          port,
		AppEnv:        appEnv,
		Version:       version,
		DatabaseURL:   databaseURL,
		WorkspaceDir:  workspaceDir,
		TemplatesDir:  templatesDir,
		JWTSecret:     jwtSecret,
		AgentToken:    agentToken,
		CatalogAppDir: catalogAppsDir,
		SecretsKey:    secretsKey,
	}
}

func (c Config) SecretsKeyBytes() ([]byte, error) {
	s := strings.TrimSpace(c.SecretsKey)
	if s == "" {
		return nil, nil
	}
	key, err := hex.DecodeString(s)
	if err != nil {
		return nil, fmt.Errorf("SAILORPORT_SECRETS_KEY must be hex: %w", err)
	}
	if len(key) != 32 {
		return nil, fmt.Errorf(
			"SAILORPORT_SECRETS_KEY must decode to 32 bytes (64 hex chars), got %d bytes",
			len(key),
		)
	}
	return key, nil
}

func (c Config) Validate() error {
	if c.AppEnv == "development" {
		return nil
	}

	var problems []string
	if c.JWTSecret == "" || c.JWTSecret == devJWTSecret {
		problems = append(problems, "AUTH_JWT_SECRET is unset or still the dev default")
	}
	if c.AgentToken == "" || c.AgentToken == devAgentToken {
		problems = append(problems, "SAILORPORT_AGENT_TOKEN is unset or still the dev default")
	}
	if strings.TrimSpace(c.SecretsKey) == "" {
		problems = append(problems, "SAILORPORT_SECRETS_KEY is unset")
	} else if _, err := c.SecretsKeyBytes(); err != nil {
		problems = append(problems, err.Error())
	}
	if len(problems) == 0 {
		return nil
	}

	return fmt.Errorf(
		"APP_ENV=%s but %s\n"+
			"hint: generate one per value with `openssl rand -hex 32`\n"+
			"hint (compose): copy deploy/compose/.env.example to deploy/compose/.env and fill it",
		c.AppEnv, strings.Join(problems, "; "),
	)
}

func defaultTemplatesDir() string {
	candidates := []string{
		"templates",
		filepath.Join("..", "..", "templates"),
		filepath.Join("..", "templates"),
	}
	for _, c := range candidates {
		if info, err := os.Stat(c); err == nil && info.IsDir() {
			if abs, err := filepath.Abs(c); err == nil {
				return abs
			}
			return c
		}
	}
	return "templates"
}

func defaultCatalogAppsDir() string {
	candidates := []string{
		"catalog-apps",
		filepath.Join("..", "..", "catalog-apps"),
		filepath.Join("..", "catalog-apps"),
	}
	for _, c := range candidates {
		if info, err := os.Stat(c); err == nil && info.IsDir() {
			if abs, err := filepath.Abs(c); err == nil {
				return abs
			}
			return c
		}
	}
	return "catalog-apps"
}

func defaultWorkspaceDir() string {
	candidates := []string{
		filepath.Join("data", "workspaces"),
		filepath.Join("..", "..", "data", "workspaces"),
		filepath.Join("..", "data", "workspaces"),
	}
	for _, c := range candidates {
		abs, err := filepath.Abs(c)
		if err != nil {
			continue
		}
		repoRoot := filepath.Dir(filepath.Dir(abs))
		if info, err := os.Stat(filepath.Join(repoRoot, "templates")); err == nil && info.IsDir() {
			return abs
		}
	}
	if abs, err := filepath.Abs(filepath.Join("..", "..", "data", "workspaces")); err == nil {
		return abs
	}
	return filepath.Join("data", "workspaces")
}

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func EnsureWorkspaceDir(dir string) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf(
			"cannot create workspace dir %q: %w\n"+
				"hint (dev): sudo chown -R \"$USER\" <repo>/data && mkdir -p <repo>/data/workspaces\n"+
				"hint (compose): use named volume sailorport_workspaces (no manual chown)",
			dir, err,
		)
	}
	probe := filepath.Join(dir, ".sailorport-write-test")
	if err := os.WriteFile(probe, []byte("ok"), 0o600); err != nil {
		return fmt.Errorf(
			"workspace dir %q is not writable: %w\n"+
				"hint (dev): sudo chown -R \"$USER\" <repo>/data\n"+
				"hint (compose): workspaces should be the named volume sailorport_workspaces",
			dir, err,
		)
	}
	_ = os.Remove(probe)
	return nil
}

func (c Config) PortInt() (int, error) {
	return strconv.Atoi(c.Port)
}
