package service

import (
	"context"
	"fmt"

	"github.com/ekkywi/sailorport/apps/api/internal/catalogapp"
)

type CatalogApps struct {
	reg *catalogapp.Registry
}

func NewCatalogApps(reg *catalogapp.Registry) *CatalogApps {
	return &CatalogApps{reg: reg}
}

func (c *CatalogApps) List(ctx context.Context) ([]catalogapp.Manifest, error) {
	_ = ctx
	return c.reg.List()
}

func (c *CatalogApps) Get(ctx context.Context, id string) (catalogapp.Manifest, error) {
	_ = ctx
	m, err := c.reg.Get(id)
	if err != nil {
		return catalogapp.Manifest{}, fmt.Errorf("%w: %v", ErrNotFound, err)
	}
	return m, nil
}
