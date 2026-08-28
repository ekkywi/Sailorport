package secrets

import (
	"context"

	"github.com/ekkywi/sailorport/apps/api/internal/model"
)

type Store interface {
	ReplaceAll(ctx context.Context, serviceID string, entries []model.CatalogEnv) error
	ResolveForDeploy(ctx context.Context, serviceID string) (map[string]string, error)
	PublicView(ctx context.Context, serviceID string) (model.CatalogEnvPublic, error)
	DeleteByServiceID(ctx context.Context, serviceID string) error
}