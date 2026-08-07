package service

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/ekkywi/sailorport/apps/api/internal/model"
	"github.com/ekkywi/sailorport/apps/api/internal/store"
)

type Workers struct {
	store *store.WorkersStore
}

func NewWorkers(s *store.WorkersStore) *Workers {
	return &Workers{store: s}
}

func (w *Workers) Register(ctx context.Context, req model.RegisterWorkerRequest) (model.Worker, error) {
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return model.Worker{}, ErrInvalid
	}

	hostname := strings.TrimSpace(req.Hostname)
	return w.store.UpsertByName(ctx, name, hostname, req.Labels)
}

func (w *Workers) Heartbeat(ctx context.Context, id, status string) (model.Worker, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return model.Worker{}, ErrInvalid
	}

	status = strings.TrimSpace(status)
	if status == "" {
		status = "online"
	}
	if status != "online" && status != "offline" && status != "draining" {
		return model.Worker{}, ErrInvalid
	}
	out, err := w.store.Heartbeat(ctx, id, status)
	if errors.Is(err, sql.ErrNoRows) {
		return model.Worker{}, ErrNotFound
	}
	return out, err
}

func (w *Workers) List(ctx context.Context) ([]model.Worker, error) {
	return w.store.List(ctx)
}
