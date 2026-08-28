package secrets

import (
	"context"
	"fmt"
	"testing"

	"github.com/ekkywi/sailorport/apps/api/internal/model"
)

type fakeCatalogEnv struct {
	byService map[string][]model.CatalogEnv
}

func (f *fakeCatalogEnv) ListByServiceID(_ context.Context, serviceID string) ([]model.CatalogEnv, error) {
	return f.byService[serviceID], nil
}

func (f *fakeCatalogEnv) ReplaceAll(context.Context, string, []model.CatalogEnv) error {
	return fmt.Errorf("not implemented")
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
