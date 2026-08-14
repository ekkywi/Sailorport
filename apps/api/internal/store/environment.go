package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ekkywi/sailorport/apps/api/internal/model"
)

type EnvironmentsStore struct {
	db *sql.DB
}

func NewEnvironmentsStore(db *sql.DB) *EnvironmentsStore {
	return &EnvironmentsStore{db: db}
}

func (s *EnvironmentsStore) List(ctx context.Context) ([]model.Environment, error) {
	const q = `
		SELECT id, slug, name, sort_order, created_at, updated_at
		FROM environments
		ORDER BY sort_order ASC, slug ASC`
	rows, err := s.db.QueryContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("List environments: %w", err)
	}

	defer rows.Close()
	out := make([]model.Environment, 0)
	for rows.Next() {
		e, err := scanEnvironment(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, e)
	}
	return out, rows.Err()
}

func scanEnvironment(row interface {
	Scan(dest ...any) error
}) (model.Environment, error) {
	var e model.Environment
	err := row.Scan(
		&e.ID,
		&e.Slug,
		&e.Name,
		&e.SortOrder,
		&e.CreatedAt,
		&e.UpdatedAt,
	)
	return e, err
}

func (s *EnvironmentsStore) GetBySlug(ctx context.Context, slug string) (model.Environment, error) {
	const q = `
		SELECT id, slug, name, sort_order, created_at, updated_at
		FROM environments
		WHERE slug = $1`

	e, err := scanEnvironment(s.db.QueryRowContext(ctx, q, slug))
	if errors.Is(err, sql.ErrNoRows) {
		return model.Environment{}, ErrNotFound
	}
	if err != nil {
		return model.Environment{}, fmt.Errorf("Get environment by slug: %w", err)
	}
	return e, nil
}
