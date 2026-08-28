package store

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/ekkywi/sailorport/apps/api/internal/model"
)

type CatalogEnvStore struct {
	db *sql.DB
}

func NewCatalogEnvStore(db *sql.DB) *CatalogEnvStore {
	return &CatalogEnvStore{db: db}
}

func (s *CatalogEnvStore) ListByServiceID(ctx context.Context, serviceID string) ([]model.CatalogEnv, error) {
	serviceID = strings.TrimSpace(serviceID)
	if serviceID == "" {
		return nil, fmt.Errorf("List catalog env: service_id is required")
	}

	const q = `
		SELECT key, value, secret
		FROM service_catalog_env
		WHERE service_id = $1::uuid
		ORDER BY key ASC
	`

	rows, err := s.db.QueryContext(ctx, q, serviceID)
	if err != nil {
		return nil, fmt.Errorf("List catalog env: %w", err)
	}
	defer rows.Close()

	out := make([]model.CatalogEnv, 0)
	for rows.Next() {
		var e model.CatalogEnv
		if err := rows.Scan(&e.Key, &e.Value, &e.Secret); err != nil {
			return nil, fmt.Errorf("Scan catalog env: %w", err)
		}
		out = append(out, e)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Iterate catalog env: %w", err)
	}
	return out, nil
}

func (s *CatalogEnvStore) ReplaceAll(ctx context.Context, serviceID string, entries []model.CatalogEnv) error {
	serviceID = strings.TrimSpace(serviceID)
	if serviceID == "" {
		return fmt.Errorf("Replace catalog env: service_id is required")
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("Replace catalog env: begin tx: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	const delQ = `DELETE FROM service_catalog_env WHERE service_id = $1::uuid`
	if _, err := tx.ExecContext(ctx, delQ, serviceID); err != nil {
		return fmt.Errorf("Replace catalog env: delete: %w", err)
	}

	const insQ = `
		INSERT INTO service_catalog_env (service_id, key, value, secret, updated_at)
		VALUES ($1::uuid, $2, $3, $4, NOW())
	`

	for _, e := range entries {
		key := strings.TrimSpace(e.Key)
		if key == "" {
			return fmt.Errorf("Replace catalog env: key is required")
		}
		if _, err := tx.ExecContext(ctx, insQ, serviceID, key, e.Value, e.Secret); err != nil {
			return fmt.Errorf("Replace catalog env: insert %q: %w", key, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("Replace catalog env: commit: %w", err)
	}
	return nil
}

func (s *CatalogEnvStore) DeleteByServiceID(ctx context.Context, serviceID string) error {
	serviceID = strings.TrimSpace(serviceID)
	if serviceID == "" {
		return fmt.Errorf("Delete catalog env: service_id is required")
	}

	const q = `DELETE FROM service_catalog_env WHERE service_id = $1::uuid`
	if _, err := s.db.ExecContext(ctx, q, serviceID); err != nil {
		return fmt.Errorf("Delete catalog env: %w", err)
	}
	return nil
}
