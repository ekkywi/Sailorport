package service

import (
	"errors"
	"strings"
	"testing"

	"github.com/ekkywi/sailorport/apps/api/internal/catalogapp"
	"github.com/ekkywi/sailorport/apps/api/internal/model"
)

func postgresManifest() catalogapp.Manifest {
	return catalogapp.Manifest{
		ID: "postgres",
		Env: []catalogapp.EnvField{
			{Name: "POSTGRES_USER", Required: false, Secret: false, Default: "postgres"},
			{Name: "POSTGRES_PASSWORD", Required: true, Secret: true},
			{Name: "POSTGRES_DB", Required: false, Secret: false, Default: "postgres"},
		},
	}
}

func TestBuildCatalogEnvEntries_RequiresPassword(t *testing.T) {
	_, err := buildCatalogEnvEntries(postgresManifest(), nil)
	if !errors.Is(err, ErrInvalid) {
		t.Fatalf("expected ErrInvalid, got %v", err)
	}
	if !strings.Contains(err.Error(), "POSTGRES_PASSWORD") {
		t.Fatalf("unexpected message: %v", err)
	}
}

func TestBuildCatalogEnvEntries_UnknownKey(t *testing.T) {
	_, err := buildCatalogEnvEntries(postgresManifest(), map[string]string{
		"POSTGRES_PASSWORD": "secret",
		"FOO":               "bar",
	})
	if !errors.Is(err, ErrInvalid) {
		t.Fatalf("expected ErrInvalid, got %v", err)
	}
	if !strings.Contains(err.Error(), "FOO") {
		t.Fatalf("unexpected message: %v", err)
	}
}

func TestBuildCatalogEnvEntries_PasswordOnlyAppliesDefaults(t *testing.T) {
	got, err := buildCatalogEnvEntries(postgresManifest(), map[string]string{
		"POSTGRES_PASSWORD": "my-secret",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 3 {
		t.Fatalf("want 3 entries, got %d: %+v", len(got), got)
	}
	byKey := map[string]model.CatalogEnv{}
	for _, e := range got {
		byKey[e.Key] = e
	}
	if byKey["POSTGRES_USER"].Value != "postgres" || byKey["POSTGRES_USER"].Secret {
		t.Fatalf("POSTGRES_USER: %+v", byKey["POSTGRES_USER"])
	}
	if byKey["POSTGRES_PASSWORD"].Value != "my-secret" || !byKey["POSTGRES_PASSWORD"].Secret {
		t.Fatalf("POSTGRES_PASSWORD: %+v", byKey["POSTGRES_PASSWORD"])
	}
	if byKey["POSTGRES_DB"].Value != "postgres" {
		t.Fatalf("POSTGRES_DB: %+v", byKey["POSTGRES_DB"])
	}
}

func TestBuildCatalogEnvEntries_CustomValues(t *testing.T) {
	got, err := buildCatalogEnvEntries(postgresManifest(), map[string]string{
		"POSTGRES_USER":     "admin",
		"POSTGRES_PASSWORD": "pw",
		"POSTGRES_DB":       "appdb",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	byKey := map[string]string{}
	for _, e := range got {
		byKey[e.Key] = e.Value
	}
	if byKey["POSTGRES_USER"] != "admin" || byKey["POSTGRES_PASSWORD"] != "pw" || byKey["POSTGRES_DB"] != "appdb" {
		t.Fatalf("values: %+v", byKey)
	}
}

func TestBuildCatalogEnvEntries_TrimsWhitespace(t *testing.T) {
	got, err := buildCatalogEnvEntries(postgresManifest(), map[string]string{
		"POSTGRES_PASSWORD":  "secret",
		"  POSTGRES_USER  ": "  trimmed  ",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, e := range got {
		if e.Key == "POSTGRES_USER" && e.Value != "trimmed" {
			t.Fatalf("value not trimmed: %+v", e)
		}
	}
}

func TestBuildCatalogEnvEntries_EmptyPasswordRejected(t *testing.T) {
	_, err := buildCatalogEnvEntries(postgresManifest(), map[string]string{
		"POSTGRES_PASSWORD": "   ",
	})
	if !errors.Is(err, ErrInvalid) {
		t.Fatalf("expected ErrInvalid, got %v", err)
	}
}

func TestBuildCatalogEnvEntries_NoEnvInManifest(t *testing.T) {
	m := catalogapp.Manifest{ID: "minimal", Env: nil}

	got, err := buildCatalogEnvEntries(m, nil)
	if err != nil || len(got) != 0 {
		t.Fatalf("empty manifest + empty input: got=%v err=%v", got, err)
	}

	_, err = buildCatalogEnvEntries(m, map[string]string{"X": "1"})
	if !errors.Is(err, ErrInvalid) {
		t.Fatalf("expected ErrInvalid for unknown key on empty manifest, got %v", err)
	}
}

func TestNormalizeCatalogEnvInput_SkipsEmptyKeys(t *testing.T) {
	got := normalizeCatalogEnvInput(map[string]string{
		"":                  "ignored",
		"POSTGRES_PASSWORD": "ok",
	})
	if len(got) != 1 || got["POSTGRES_PASSWORD"] != "ok" {
		t.Fatalf("normalize: %+v", got)
	}
}

func existingPostgresEnv(password string) []model.CatalogEnv {
	return []model.CatalogEnv{
		{Key: "POSTGRES_USER", Value: "postgres", Secret: false},
		{Key: "POSTGRES_PASSWORD", Value: password, Secret: true},
		{Key: "POSTGRES_DB", Value: "postgres", Secret: false},
	}
}

func TestMergeCatalogEnvForUpdate_KeepSecretWhenEmpty(t *testing.T) {
	got, err := mergeCatalogEnvForUpdate(postgresManifest(), map[string]string{
		"POSTGRES_USER": "admin",
	}, existingPostgresEnv("old-secret"))
	if err != nil {
		t.Fatal(err)
	}
	byKey := map[string]string{}
	for _, e := range got {
		byKey[e.Key] = e.Value
	}
	if byKey["POSTGRES_USER"] != "admin" {
		t.Fatalf("user: %v", byKey)
	}
	if byKey["POSTGRES_PASSWORD"] != "old-secret" {
		t.Fatalf("password should be kept, got %q", byKey["POSTGRES_PASSWORD"])
	}
}

func TestMergeCatalogEnvForUpdate_ReplacesSecretWhenProvided(t *testing.T) {
	got, err := mergeCatalogEnvForUpdate(postgresManifest(), map[string]string{
		"POSTGRES_PASSWORD": "new-secret",
	}, existingPostgresEnv("old-secret"))
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range got {
		if e.Key == "POSTGRES_PASSWORD" && e.Value != "new-secret" {
			t.Fatalf("want new-secret, got %q", e.Value)
		}
	}
}

func TestMergeCatalogEnvForUpdate_RequiredSecretMissingWithoutExisting(t *testing.T) {
	_, err := mergeCatalogEnvForUpdate(postgresManifest(), map[string]string{
		"POSTGRES_USER": "admin",
	}, nil)
	if !errors.Is(err, ErrInvalid) {
		t.Fatalf("expected ErrInvalid, got %v", err)
	}
}
func TestMergeCatalogEnvForUpdate_UnknownKey(t *testing.T) {
	_, err := mergeCatalogEnvForUpdate(postgresManifest(), map[string]string{
		"POSTGRES_PASSWORD": "x",
		"BAD":               "y",
	}, existingPostgresEnv("old"))
	if !errors.Is(err, ErrInvalid) {
		t.Fatalf("expected ErrInvalid, got %v", err)
	}
}