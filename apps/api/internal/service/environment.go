package service

import (
	"context"

	"github.com/ekkywi/sailorport/apps/api/internal/model"
	"github.com/ekkywi/sailorport/apps/api/internal/store"
)

type Environments struct {
	store *store.EnvironmentsStore
}

func NewEnvironments(s *store.EnvironmentsStore) *Environments {
	return &Environments{store: s}
}

func (e *Environments) List(ctx context.Context) ([]model.Environment, error) {
	return e.store.List(ctx)
}
