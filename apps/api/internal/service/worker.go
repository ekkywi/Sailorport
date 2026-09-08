package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
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
		status = model.WorkerStatusOnline
	}
	if !model.IsWorkerStatus(status) {
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

func (w *Workers) Get(ctx context.Context, id string) (model.Worker, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return model.Worker{}, ErrInvalid
	}
	out, err := w.store.Get(ctx, id)
	if errors.Is(err, store.ErrNotFound) {
		return model.Worker{}, ErrNotFound
	}
	if err != nil {
		return model.Worker{}, err
	}
	return out, nil
}

func mergeWorkerCapabilityLabels(existing map[string]any, req model.UpdateWorkerLabelsRequest) map[string]any {
	out := make(map[string]any, len(existing)+2)
	for k, v := range existing {
		out[k] = v
	}
	out["tier"] = strings.TrimSpace(req.Tier)
	out["environments"] = strings.TrimSpace(req.Environments)
	return out
}

func (w *Workers) UpdateLabels(ctx context.Context, id string, req model.UpdateWorkerLabelsRequest) (model.Worker, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return model.Worker{}, ErrInvalid
	}

	existing, err := w.Get(ctx, id)
	if err != nil {
		return model.Worker{}, err
	}

	merged := mergeWorkerCapabilityLabels(existing.Labels, req)
	out, err := w.store.UpdateLabels(ctx, id, merged)
	if errors.Is(err, store.ErrNotFound) {
		return model.Worker{}, ErrNotFound
	}
	return out, err
}

func (w *Workers) Decommission(ctx context.Context, id string) (model.Worker, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return model.Worker{}, ErrInvalid
	}
	if _, err := w.Get(ctx, id); err != nil {
		return model.Worker{}, err
	}
	out, err := w.store.SetStatus(ctx, id, model.WorkerStatusDraining)
	if errors.Is(err, store.ErrNotFound) {
		return model.Worker{}, ErrNotFound
	}
	return out, err
}

func (w *Workers) Restore(ctx context.Context, id string) (model.Worker, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return model.Worker{}, ErrInvalid
	}
	existing, err := w.Get(ctx, id)
	if err != nil {
		return model.Worker{}, err
	}
	if existing.Status != model.WorkerStatusDraining {
		return model.Worker{}, fmt.Errorf(
			"%w: worker %q is not draining (status=%s)",
			ErrInvalid, existing.Name, existing.Status,
		)
	}
	out, err := w.store.SetStatus(ctx, id, model.WorkerStatusOffline)
	if errors.Is(err, store.ErrNotFound) {
		return model.Worker{}, ErrNotFound
	}
	return out, err
}

func workerDeployConflict(w model.Worker, requireOnline bool) error {
	if w.Status == model.WorkerStatusDraining {
		return fmt.Errorf("%w: worker %q is decommissioned (draining)", ErrConflict, w.Name)
	}
	if requireOnline && !model.WorkerAcceptsDeploy(w.Status) {
		return fmt.Errorf("%w: worker %q is %s (must be online)", ErrConflict, w.Name, w.Status)
	}
	return nil
}
