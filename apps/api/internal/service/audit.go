package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/ekkywi/sailorport/apps/api/internal/model"
	"github.com/ekkywi/sailorport/apps/api/internal/store"
)

type Audit struct {
	store *store.AuditStore
}

func NewAudit(s *store.AuditStore) *Audit {
	return &Audit{store: s}
}

func (a *Audit) Record(ctx context.Context, rec model.AuditRecord) error {
	if a == nil || a.store == nil {
		return nil
	}
	rec.Action = strings.TrimSpace(rec.Action)
	rec.ResourceType = strings.TrimSpace(rec.ResourceType)
	if rec.Action == "" || rec.ResourceType == "" {
		return fmt.Errorf("%w: audit action and resource_type are required", ErrInvalid)
	}
	if rec.Payload == nil {
		rec.Payload = map[string]any{}
	}
	_, err := a.store.Insert(ctx, rec)
	if err != nil {
		return fmt.Errorf("record audit: %w", err)
	}
	return nil
}

func (a *Audit) List(ctx context.Context, limit int) ([]model.AuditEvent, error) {
	if a == nil || a.store == nil {
		return []model.AuditEvent{}, nil
	}
	out, err := a.store.List(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("list audit: %w", err)
	}
	return out, nil
}
