package secrets

import (
	"context"
	"testing"

	"github.com/ekkywi/sailorport/apps/api/internal/model"
)

type fakeCatalogEnv struct {
	byService map[string][]model.CatalogEnv
}

func (f *fakeCatalogEnv) ListByServiceID(_ context.Context, serviceID string) ([]model.CatalogEnv, error) {
	return f.byService[serviceID], nil
}

func (f *fakeCatalogEnv) ReplaceAll(_ context.Context, serviceID string, entries []model.CatalogEnv) error {
	copied := make([]model.CatalogEnv, len(entries))
	copy(copied, entries)
	f.byService[serviceID] = copied
	return nil
}

func (f *fakeCatalogEnv) DeleteByServiceID(context.Context, string) error {
	return nil
}

func TestCatalogEnvPublic_RedactsSecrets(t *testing.T) {
	pub := catalogEnvPublic([]model.CatalogEnv{
		{Key: "POSTGRES_USER", Value: "postgres", Secret: false},
		{Key: "POSTGRES_PASSWORD", Value: "secret", Secret: true},
	})
	if pub["POSTGRES_USER"] != "postgres" {
		t.Fatalf("user: %v", pub["POSTGRES_USER"])
	}
	if pub["POSTGRES_PASSWORD_set"] != true {
		t.Fatal("expected POSTGRES_PASSWORD_set")
	}
	if _, leak := pub["POSTGRES_PASSWORD"]; leak {
		t.Fatal("secret value must not appear")
	}
}

func TestResolveForDeploy_IncludesSecrets(t *testing.T) {
	fake := &fakeCatalogEnv{byService: map[string][]model.CatalogEnv{
		"svc-1": {{Key: "POSTGRES_PASSWORD", Value: "secret", Secret: true}},
	}}
	p := &PlaintextStore{env: fake}
	got, err := p.ResolveForDeploy(context.Background(), "svc-1")
	if err != nil || got["POSTGRES_PASSWORD"] != "secret" {
		t.Fatalf("ResolveForDeploy: %v %v", got, err)
	}
}

func TestEncryptedStore_ReplaceAll_EncryptsSecretsOnly(t *testing.T) {
	fake := &fakeCatalogEnv{byService: map[string][]model.CatalogEnv{}}
	store, err := NewEncrypted(fake, testSecretsKey(t))
	if err != nil {
		t.Fatal(err)
	}

	err = store.ReplaceAll(context.Background(), "svc-1", []model.CatalogEnv{
		{Key: "POSTGRES_USER", Value: "postgres", Secret: false},
		{Key: "POSTGRES_PASSWORD", Value: "secret123", Secret: true},
	})
	if err != nil {
		t.Fatal(err)
	}

	stored := fake.byService["svc-1"]
	if len(stored) != 2 {
		t.Fatalf("stored %d rows", len(stored))
	}
	byKey := map[string]model.CatalogEnv{}
	for _, e := range stored {
		byKey[e.Key] = e
	}
	if byKey["POSTGRES_USER"].Value != "postgres" {
		t.Fatalf("non-secret should stay plaintext, got %q", byKey["POSTGRES_USER"].Value)
	}
	if !isEncryptedValue(byKey["POSTGRES_PASSWORD"].Value) {
		t.Fatalf("secret should be encrypted, got %q", byKey["POSTGRES_PASSWORD"].Value)
	}
}

func TestEncryptedStore_ResolveForDeploy_Decrypts(t *testing.T) {
	fake := &fakeCatalogEnv{byService: map[string][]model.CatalogEnv{}}
	store, err := NewEncrypted(fake, testSecretsKey(t))
	if err != nil {
		t.Fatal(err)
	}

	if err := store.ReplaceAll(context.Background(), "svc-1", []model.CatalogEnv{
		{Key: "POSTGRES_PASSWORD", Value: "secret123", Secret: true},
	}); err != nil {
		t.Fatal(err)
	}

	got, err := store.ResolveForDeploy(context.Background(), "svc-1")
	if err != nil {
		t.Fatal(err)
	}
	if got["POSTGRES_PASSWORD"] != "secret123" {
		t.Fatalf("deploy map: %v", got)
	}
}

func TestEncryptedStore_PublicView_RedactsSecrets(t *testing.T) {
	fake := &fakeCatalogEnv{byService: map[string][]model.CatalogEnv{}}
	store, err := NewEncrypted(fake, testSecretsKey(t))
	if err != nil {
		t.Fatal(err)
	}

	if err := store.ReplaceAll(context.Background(), "svc-1", []model.CatalogEnv{
		{Key: "POSTGRES_USER", Value: "postgres", Secret: false},
		{Key: "POSTGRES_PASSWORD", Value: "secret123", Secret: true},
	}); err != nil {
		t.Fatal(err)
	}

	pub, err := store.PublicView(context.Background(), "svc-1")
	if err != nil {
		t.Fatal(err)
	}
	if pub["POSTGRES_USER"] != "postgres" {
		t.Fatalf("user: %v", pub["POSTGRES_USER"])
	}
	if pub["POSTGRES_PASSWORD_set"] != true {
		t.Fatal("expected POSTGRES_PASSWORD_set")
	}
}

func TestEncryptedStore_BackwardCompat_PlaintextInDB(t *testing.T) {
	fake := &fakeCatalogEnv{byService: map[string][]model.CatalogEnv{
		"svc-old": {
			{Key: "POSTGRES_PASSWORD", Value: "old-plain", Secret: true},
		},
	}}
	store, err := NewEncrypted(fake, testSecretsKey(t))
	if err != nil {
		t.Fatal(err)
	}

	got, err := store.ResolveForDeploy(context.Background(), "svc-old")
	if err != nil {
		t.Fatal(err)
	}
	if got["POSTGRES_PASSWORD"] != "old-plain" {
		t.Fatalf("legacy plaintext: %v", got)
	}
}
