package config

import (
	"os"
	"path/filepath"
	"strconv"
)

type Config struct {
	Port         string
	AppEnv       string
	Version      string
	DatabaseURL  string
	WorkspaceDir string
	TemplatesDir string
	JWTSecret    string
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
	jwtSecret := getenv("AUTH_JWT_SECRET", "dev-only-change-me")

	return Config{
		Port:         port,
		AppEnv:       appEnv,
		Version:      version,
		DatabaseURL:  databaseURL,
		WorkspaceDir: workspaceDir,
		TemplatesDir: templatesDir,
		JWTSecret:    jwtSecret,
	}
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

func (c Config) PortInt() (int, error) {
	return strconv.Atoi(c.Port)
}
