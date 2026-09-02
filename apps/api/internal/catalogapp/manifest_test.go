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
