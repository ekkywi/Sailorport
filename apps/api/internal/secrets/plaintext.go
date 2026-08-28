package secrets

import (
	"context"

	"github.com/ekkywi/sailorport/apps/api/internal/model"
	"github.com/ekkywi/sailorport/apps/api/internal/store"
)

type PlaintextStore struct {
	env catalogEnvStore
}

func NewPlaintext(env *store.CatalogEnvStore) *PlaintextStore {
	return &PlaintextStore{env: env}
}

func (p *PlaintextStore) ReplaceAll(ctx context.Context, serviceID string, entries []model.CatalogEnv) error {
	return p.env.ReplaceAll(ctx, serviceID, entries)
}

func (p *PlaintextStore) DeleteByServiceID(ctx context.Context, serviceID string) error {
	return p.env.DeleteByServiceID(ctx, serviceID)
}

func (p *PlaintextStore) ResolveForDeploy(ctx context.Context, serviceID string) (map[string]string, error) {
	rows, err := p.env.ListByServiceID(ctx, serviceID)
	if err != nil {
		return nil, err
	}
	out := make(map[string]string, len(rows))
	for _, e := range rows {
		out[e.Key] = e.Value
	}
	return out, nil
}

func (p *PlaintextStore) PublicView(ctx context.Context, serviceID string) (model.CatalogEnvPublic, error) {
	rows, err := p.env.ListByServiceID(ctx, serviceID)
	if err != nil {
		return nil, err
	}
	return catalogEnvPublic(rows), nil
}

func catalogEnvPublic(rows []model.CatalogEnv) model.CatalogEnvPublic {
	if len(rows) == 0 {
		return nil
	}
	out := make(model.CatalogEnvPublic, len(rows)*2)
	for _, e := range rows {
		if e.Secret {
			out[e.Key+"_set"] = true
			continue
		}
		out[e.Key] = e.Value
	}
	return out
}

var _ Store = (*PlaintextStore)(nil)
