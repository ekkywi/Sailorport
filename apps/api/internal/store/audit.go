package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/ekkywi/sailorport/apps/api/internal/model"
)

type AuditStore struct {
	db *sql.DB
}

func NewAuditStore(db *sql.DB) *AuditStore {
	return &AuditStore{db: db}
}

func (s *AuditStore) Insert(ctx context.Context, rec model.AuditRecord) (model.AuditEvent, error) {
	payload := rec.Payload
	if payload == nil {
		payload = map[string]any{}
	}
	raw, err := json.Marshal(payload)
	if err != nil {
		return model.AuditEvent{}, fmt.Errorf("Insert audit: marshal payload: %w", err)
	}

	const q = `
		INSERT INTO audit_events (
			actor_id, actor_email, action, resource_type, resource_id, resource_name, payload
		)
		VALUES (
			$1, $2, $3, $4, $5, $6, $7
		)
		RETURNING id, at, actor_id, actor_email, action, resource_type, resource_id, resource_name, payload
	`

	var actorID any
	if rec.ActorID != "" {
		actorID = rec.ActorID
	} else {
		actorID = nil
	}

	row := s.db.QueryRowContext(
		ctx,
		q,
		actorID,
		rec.ActorEmail,
		rec.Action,
		rec.ResourceType,
		rec.ResourceID,
		rec.ResourceName,
		raw,
	)
	return scanAuditEvent(row)
}

func scanAuditEvent(row interface{ Scan(dest ...any) error }) (model.AuditEvent, error) {
	var (
		ev      model.AuditEvent
		actorID sql.NullString
		raw     []byte
	)
	err := row.Scan(
		&ev.ID,
		&ev.At,
		&actorID,
		&ev.ActorEmail,
		&ev.Action,
		&ev.ResourceType,
		&ev.ResourceID,
		&ev.ResourceName,
		&raw,
	)
	if err != nil {
		return model.AuditEvent{}, err
	}
	if actorID.Valid {
		ev.ActorID = actorID.String
	}
	if len(raw) > 0 {
		if err := json.Unmarshal(raw, &ev.Payload); err != nil {
			return model.AuditEvent{}, fmt.Errorf("scan audit payload: %w", err)
		}
	}
	if ev.Payload == nil {
		ev.Payload = map[string]any{}
	}
	return ev, nil
}

func (s *AuditStore) List(ctx context.Context, limit int) ([]model.AuditEvent, error) {
	if limit <= 0 {
		limit = 50
	}
	if limit > 200 {
		limit = 200
	}

	const q = `
		SELECT id, at, actor_id, actor_email, action, resource_type, resource_id, resource_name, payload
		FROM audit_events
		ORDER BY at DESC
		LIMIT $1
	`

	rows, err := s.db.QueryContext(ctx, q, limit)
	if err != nil {
		return nil, fmt.Errorf("List audit events: %w", err)
	}
	defer rows.Close()

	out := make([]model.AuditEvent, 0)
	for rows.Next() {
		ev, err := scanAuditEvent(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, ev)
	}
	return out, rows.Err()
}