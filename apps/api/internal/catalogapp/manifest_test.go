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
	// go test cwd is the package dir (apps/api/internal/catalogapp)
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
