package secrets

import (
	"context"
	"testing"

	"github.com/ekkywi/sailorport/apps/api/internal/model"
)

// Step 25d smoke: encrypted at-rest round-trip (secret in DB ciphertext, deploy/API behavior unchanged).
func TestSmoke_EncryptedCatalogEnvAtRest(t *testing.T) {
	fake := &fakeCatalogEnv{byService: map[string][]model.CatalogEnv{}}
	store, err := NewEncrypted(fake, testSecretsKey(t))
	if err != nil {
		t.Fatal(err)
	}

	const password = "encrypt-smoke-123"
	if err := store.ReplaceAll(context.Background(), "svc-smoke", []model.CatalogEnv{
		{Key: "POSTGRES_USER", Value: "postgres", Secret: false},
		{Key: "POSTGRES_PASSWORD", Value: password, Secret: true},
	}); err != nil {
		t.Fatal(err)
	}

	stored := fake.byService["svc-smoke"]
	byKey := map[string]model.CatalogEnv{}
	for _, e := range stored {
		byKey[e.Key] = e
	}
	if byKey["POSTGRES_USER"].Value != "postgres" {
		t.Fatalf("non-secret must stay plaintext in DB, got %q", byKey["POSTGRES_USER"].Value)
	}
	if !isEncryptedValue(byKey["POSTGRES_PASSWORD"].Value) {
		t.Fatalf("secret must be encrypted in DB, got %q", byKey["POSTGRES_PASSWORD"].Value)
	}

	deploy, err := store.ResolveForDeploy(context.Background(), "svc-smoke")
	if err != nil {
		t.Fatal(err)
	}
	if deploy["POSTGRES_PASSWORD"] != password {
		t.Fatalf("deploy must get plaintext password, got %q", deploy["POSTGRES_PASSWORD"])
	}

	pub, err := store.PublicView(context.Background(), "svc-smoke")
	if err != nil {
		t.Fatal(err)
	}
	if pub["POSTGRES_USER"] != "postgres" {
		t.Fatalf("public user: %v", pub["POSTGRES_USER"])
	}
	if pub["POSTGRES_PASSWORD_set"] != true {
		t.Fatal("public view must redact secret")
	}
	if _, leak := pub["POSTGRES_PASSWORD"]; leak {
		t.Fatal("public view must not leak secret value")
	}
}
