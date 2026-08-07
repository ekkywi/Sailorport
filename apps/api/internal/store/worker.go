package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/ekkywi/sailorport/apps/api/internal/model"
)

type WorkersStore struct {
	db *sql.DB
}

func NewWorkersStore(db *sql.DB) *WorkersStore {
	return &WorkersStore{db: db}
}

func (s *WorkersStore) UpsertByName(ctx context.Context, name, hostname string, labels map[string]any) (model.Worker, error) {
	if labels == nil {
		labels = map[string]any{}
	}
	raw, err := json.Marshal(labels)
	if err != nil {
		return model.Worker{}, err
	}

	const q = `
		INSERT INTO workers (name, hostname, labels, status, updated_at)
		VALUES ($1, $2, $3::jsonb, 'offline', NOW())
		ON CONFLICT (name) DO UPDATE SET
		hostname = EXCLUDED.hostname,
		labels = EXCLUDED.labels,
		updated_at = NOW()
		RETURNING id, name, hostname, status, labels, last_seen_at, created_at, updated_at
		`
	return scanWorker(s.db.QueryRowContext(ctx, q, name, hostname, raw))
}

func (s *WorkersStore) Heartbeat(ctx context.Context, id, status string) (model.Worker, error) {
	const q = `
		UPDATE workers
		SET status = $2, last_seen_at = NOW(), updated_at = NOW()
		WHERE id = $1
		RETURNING id, name, hostname, status, labels, last_seen_at, created_at, updated_at
		`
	w, err := scanWorker(s.db.QueryRowContext(ctx, q, id, status))
	if errors.Is(err, sql.ErrNoRows) {
		return model.Worker{}, sql.ErrNoRows
	}
	return w, err
}

func (s *WorkersStore) List(ctx context.Context) ([]model.Worker, error) {
	const q = `
		SELECT id, name, hostname, status, labels, last_seen_at, created_at, updated_at
		FROM workers
		ORDER BY created_at DESC
		`
	rows, err := s.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]model.Worker, 0)
	for rows.Next() {
		w, err := scanWorker(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, w)
	}
	return out, rows.Err()
}

type rowScanner interface {
	Scan(dest ...any) error
}

func scanWorker(row rowScanner) (model.Worker, error) {
	var (
		w        model.Worker
		raw      []byte
		lastSeen sql.NullTime
	)
	err := row.Scan(
		&w.ID, &w.Name, &w.Hostname, &w.Status, &raw, &lastSeen, &w.CreatedAt, &w.UpdatedAt,
	)
	if err != nil {
		return model.Worker{}, err
	}
	if lastSeen.Valid {
		t := lastSeen.Time
		w.LastSeenAt = &t
	}
	w.Labels = map[string]any{}
	if len(raw) > 0 {
		_ = json.Unmarshal(raw, &w.Labels)
	}
	return w, nil
}
