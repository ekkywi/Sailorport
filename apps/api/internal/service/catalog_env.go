package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/ekkywi/sailorport/apps/api/internal/catalogapp"
	"github.com/ekkywi/sailorport/apps/api/internal/model"
)

func buildCatalogEnvEntries(m catalogapp.Manifest, input map[string]string) ([]model.CatalogEnv, error) {
	normalized := normalizeCatalogEnvInput(input)

	if len(m.Env) == 0 {
		if len(normalized) > 0 {
			for k := range normalized {
				return nil, fmt.Errorf("%w: unknown catalog env key %q", ErrInvalid, k)
			}
		}
		return nil, nil
	}

	allowed := make(map[string]catalogapp.EnvField, len(m.Env))
	for _, f := range m.Env {
		allowed[f.Name] = f
	}

	out := make([]model.CatalogEnv, 0, len(m.Env))
	for _, f := range m.Env {
		val := normalized[f.Name]
		if val == "" {
			val = f.Default
		}
		if f.Required && val == "" {
			return nil, fmt.Errorf("%w: catalog env %q is required", ErrInvalid, f.Name)
		}
		if val == "" {
			continue
		}
		out = append(out, model.CatalogEnv{
			Key:    f.Name,
			Value:  val,
			Secret: f.Secret,
		})
	}

	for k := range normalized {
		if _, ok := allowed[k]; !ok {
			return nil, fmt.Errorf("%w: unknown catalog env key %q", ErrInvalid, k)
		}
	}

	return out, nil
}

func mergeCatalogEnvForUpdate(
	m catalogapp.Manifest,
	input map[string]string,
	existing []model.CatalogEnv,
) ([]model.CatalogEnv, error) {
	normalized := normalizeCatalogEnvInput(input)

	if len(m.Env) == 0 {
		if len(normalized) > 0 {
			for k := range normalized {
				return nil, fmt.Errorf("%w: unknown catalog env key %q", ErrInvalid, k)
			}
		}
		return nil, nil
	}

	existingByKey := make(map[string]string, len(existing))
	for _, e := range existing {
		existingByKey[e.Key] = e.Value
	}

	allowed := make(map[string]catalogapp.EnvField, len(m.Env))
	for _, f := range m.Env {
		allowed[f.Name] = f
	}

	out := make([]model.CatalogEnv, 0, len(m.Env))
	for _, f := range m.Env {
		val := normalized[f.Name]

		if val == "" && f.Secret {
			if prev, ok := existingByKey[f.Name]; ok && prev != "" {
				val = prev
			}
		}
		if val == "" {
			val = f.Default
		}
		if f.Required && val == "" {
			return nil, fmt.Errorf("%w: catalog env %q is required", ErrInvalid, f.Name)
		}
		if val == "" {
			continue
		}
		out = append(out, model.CatalogEnv{
			Key: f.Name,
			Value: val,
			Secret: f.Secret,
		})
	}

	for k := range normalized {
		if _, ok := allowed[k]; !ok {
			return nil, fmt.Errorf("%w: unknown catalog env key %q", ErrInvalid, k)
		}
	}
	
	return out, nil
}

func normalizeCatalogEnvInput(input map[string]string) map[string]string {
	if len(input) == 0 {
		return map[string]string{}
	}
	out := make(map[string]string, len(input))
	for k, v := range input {
		key := strings.TrimSpace(k)
		if key == "" {
			continue
		}
		out[key] = strings.TrimSpace(v)
	}
	return out
}

func (c *Catalog) attachCatalogEnvPublic(ctx context.Context, svc *model.Service) error {
	if c.secrets == nil || svc.SourceType != "catalog_app" {
		return nil
	}
	pub, err := c.secrets.PublicView(ctx, svc.ID)
	if err != nil {
		return fmt.Errorf("catalog env public view: %w", err)
	}
	svc.CatalogEnv = pub
	return nil
}