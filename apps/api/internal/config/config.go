package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Nilai bawaan ini ada di repo publik, jadi hanya boleh dipakai di development.
// Validate() menolaknya untuk APP_ENV lain.
const (
	devJWTSecret  = "dev-only-change-me"
	devAgentToken = "dev-agent-token"
)

type Config struct {
	Port         string
	AppEnv       string
	Version      string
	DatabaseURL  string
	WorkspaceDir string
	TemplatesDir string
	JWTSecret    string
	AgentToken   string
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

	return Config{
		Port:         port,
		AppEnv:       appEnv,
		Version:      version,
		DatabaseURL:  databaseURL,
		WorkspaceDir: workspaceDir,
		TemplatesDir: templatesDir,
		JWTSecret:    jwtSecret,
		AgentToken:   agentToken,
	}
}

// Validate menolak secret bawaan di luar development. Tanpa ini, siapa pun yang
// membaca repo bisa menandatangani JWT role=admin atau memakai agent token.
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

// EnsureWorkspaceDir creates the workspace root and verifies it is writable.
// Call at process start so scaffold fails fast with a clear message.
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
