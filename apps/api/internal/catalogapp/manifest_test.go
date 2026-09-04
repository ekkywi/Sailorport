package catalogapp

import (
	"strings"
	"testing"
)

func TestValidateEnvFields_EmptyOK(t *testing.T) {
	if err := validateEnvFields("postgres", nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateEnvFields_DuplicateName(t *testing.T) {
	err := validateEnvFields("postgres", []EnvField{
		{Name: "POSTGRES_PASSWORD"},
		{Name: "POSTGRES_PASSWORD"},
	})
	if err == nil {
		t.Fatal("expected duplicate error")
	}
	if !strings.Contains(err.Error(), "duplicate") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateEnvFields_InvalidName(t *testing.T) {
	err := validateEnvFields("postgres", []EnvField{
		{Name: "bad-name"},
	})
	if err == nil {
		t.Fatal("expected invalid name error")
	}
}

func TestRegistry_GetPostgresManifestEnv(t *testing.T) {
	root := findCatalogAppsRoot(t)
	reg := NewRegistry(root)
	m, err := reg.Get("postgres")
	if err != nil {
		t.Fatalf("Get postgres: %v", err)
	}
	if len(m.Env) < 2 {
		t.Fatalf("expected env schema, got %d fields", len(m.Env))
	}
	var hasPassword bool
	for _, f := range m.Env {
		if f.Name == "POSTGRES_PASSWORD" {
			hasPassword = true
			if !f.Required || !f.Secret {
				t.Fatalf("POSTGRES_PASSWORD: want required+secret, got required=%v secret=%v", f.Required, f.Secret)
			}
		}
	}
	if !hasPassword {
		t.Fatal("missing POSTGRES_PASSWORD in manifest env")
	}
}

func findCatalogAppsRoot(t *testing.T) string {
	t.Helper()
	candidates := []string{
		"../../../../catalog-apps",
		"../../../catalog-apps",
		"../../catalog-apps",
		"catalog-apps",
	}
	for _, p := range candidates {
		if _, err := NewRegistry(p).Get("postgres"); err == nil {
			return p
		}
	}
	t.Fatal("could not find catalog-apps/postgres from test cwd")
	return ""
}

func TestValidateVersions_EmptyOK(t *testing.T) {
	if err := validateVersions("postgres", "postgres:16-alpine", nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateVersions_RequiresOneDefault(t *testing.T) {
	err := validateVersions("postgres", "postgres:16-alpine", []Version{
		{Tag: "16-alpine", Image: "postgres:16-alpine"},
	})
	if err == nil || !strings.Contains(err.Error(), "exactly one default") {
		t.Fatalf("expected default error, got %v", err)
	}
}

func TestValidateVersions_DuplicateTag(t *testing.T) {
	err := validateVersions("postgres", "postgres:16-alpine", []Version{
		{Tag: "16-alpine", Image: "postgres:16-alpine", Default: true},
		{Tag: "16-alpine", Image: "postgres:16", Default: false},
	})
	if err == nil || !strings.Contains(err.Error(), "duplicate") {
		t.Fatalf("expected duplicate error, got %v", err)
	}
}

func TestValidateVersions_TopImageMustMatchDefault(t *testing.T) {
	err := validateVersions("postgres", "postgres:15-alpine", []Version{
		{Tag: "16-alpine", Image: "postgres:16-alpine", Default: true},
	})
	if err == nil || !strings.Contains(err.Error(), "top-level image") {
		t.Fatalf("expected top-level mismatch error, got %v", err)
	}
}

func TestValidateVersions_OK(t *testing.T) {
	err := validateVersions("postgres", "postgres:16-alpine", []Version{
		{Tag: "15-alpine", Image: "postgres:15-alpine"},
		{Tag: "16-alpine", Image: "postgres:16-alpine", Default: true},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidteCommand_EmptyOK(t *testing.T) {
	if err := validateCommand("postgres", nil, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateCommand_LiteralOk(t *testing.T) {
	err := validateCommand("redis", []EnvField{{Name: "REDIS_PASSWORD"}}, []string{
		"redis-server",
		"--requirepass",
		"${REDIS_PASSWORD}",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateCommand_UnkownPlaceholder(t *testing.T) {
	err := validateCommand("redis", []EnvField{{Name: "REDIS_PASSWORD"}}, []string{
		"${UNKNOWN}",
	})
	if err == nil || !strings.Contains(err.Error(), "not defined in env") {
		t.Fatalf("expected undefined placeholder error, got %v", err)
	}
}

func TestValidateCommand_Unclosed(t *testing.T) {
	err := validateCommand("redis", []EnvField{{Name: "REDIS_PASSWORD"}}, []string{
		"${REDIS_PASSWORD",
	})
	if err == nil || !strings.Contains(err.Error(), "unclosed") {
		t.Fatalf("expected unclosed error, got %v", err)
	}
}

func TestValidateCommand_EmptyArg(t *testing.T) {
	err := validateCommand("redis", nil, []string{"redis-server", " "})
	if err == nil || !strings.Contains(err.Error(), "empty argument") {
		t.Fatalf("expected empty argument error, got %v", err)
	}
}

func TestRegistry_GetPostgresStillOK(t *testing.T) {
	root := findCatalogAppsRoot(t)
	m, err := NewRegistry(root).Get("postgres")
	if err != nil {
		t.Fatalf("postgres must remain valid without command: %v", err)
	}
	if len(m.Command) != 0 {
		t.Fatalf("postgres should have no command, got %#v", m.Command)
	}
}

func TestRegistry_GetRedisManifest(t *testing.T) {
	root := findCatalogAppsRoot(t)
	m, err := NewRegistry(root).Get("redis")
	if err != nil {
		t.Fatalf("Get redis: %v", err)
	}
	if m.ContainerPort != 6379 {
		t.Fatalf("port: want 6379, got %d", m.ContainerPort)
	}
	if len(m.Command) < 3 {
		t.Fatalf("expected command with requirepass, got %#v", m.Command)
	}
	var hasPass bool
	for _, f := range m.Env {
		if f.Name == "REDIS_PASSWORD" && f.Secret && f.Required {
			hasPass = true
		}
	}
	if !hasPass {
		t.Fatal("missing REDIS_PASSWORD secret required")
	}
}
