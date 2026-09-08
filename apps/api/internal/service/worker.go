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

// Allowed admin-editable worker tiers (empty = unset).
var allowedWorkerTiers = map[string]struct{}{
	"":        {},
	"nonprod": {},
	"prod":    {},
}

type Workers struct {
	store *store.WorkersStore
	envs  *store.EnvironmentsStore
}

func NewWorkers(s *store.WorkersStore, envs *store.EnvironmentsStore) *Workers {
	return &Workers{store: s, envs: envs}
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

// normalizeWorkerEnvironmentsInput parses a comma-separated env list into
// unique lowercase slugs (preserving first-seen order).
func normalizeWorkerEnvironmentsInput(raw string) []string {
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	seen := map[string]struct{}{}
	for _, p := range parts {
		slug := strings.ToLower(strings.TrimSpace(p))
		if slug == "" {
			continue
		}
		if _, ok := seen[slug]; ok {
			continue
		}
		seen[slug] = struct{}{}
		out = append(out, slug)
	}
	return out
}

func validateWorkerTier(tier string) error {
	tier = strings.TrimSpace(tier)
	if _, ok := allowedWorkerTiers[tier]; !ok {
		return fmt.Errorf("%w: tier must be empty, \"nonprod\", or \"prod\"", ErrInvalid)
	}
	return nil
}

func (w *Workers) validateWorkerEnvironments(ctx context.Context, raw string) (string, error) {
	slugs := normalizeWorkerEnvironmentsInput(raw)
	if len(slugs) == 0 {
		return "", nil // empty = allow all environments
	}
	if w.envs == nil {
		return "", fmt.Errorf("environments store not configured")
	}
	envs, err := w.envs.List(ctx)
	if err != nil {
		return "", err
	}
	known := make(map[string]struct{}, len(envs))
	for _, e := range envs {
		known[strings.ToLower(strings.TrimSpace(e.Slug))] = struct{}{}
	}
	for _, slug := range slugs {
		if _, ok := known[slug]; !ok {
			return "", fmt.Errorf("%w: unknown environment %q", ErrInvalid, slug)
		}
	}
	return strings.Join(slugs, ","), nil
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

	req.Tier = strings.TrimSpace(req.Tier)
	if err := validateWorkerTier(req.Tier); err != nil {
		return model.Worker{}, err
	}
	normalizedEnvs, err := w.validateWorkerEnvironments(ctx, req.Environments)
	if err != nil {
		return model.Worker{}, err
	}
	req.Environments = normalizedEnvs

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
