package secrets

import (
	"context"

	"github.com/ekkywi/sailorport/apps/api/internal/model"
)

type catalogEnvStore interface {
	ListByServiceID(ctx context.Context, serviceID string) ([]model.CatalogEnv, error)
	ReplaceAll(ctx context.Context, serviceID string, entries []model.CatalogEnv) error
	DeleteByServiceID(ctx context.Context, serviceID string) error
}