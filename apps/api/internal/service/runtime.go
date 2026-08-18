package service

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/ekkywi/sailorport/apps/api/internal/model"
	"github.com/ekkywi/sailorport/apps/api/internal/store"
)

// latestPerEnv picks the newest deployment per environment slug.
// deps must be ordered by created_at DESC (ListByService already does this).
func latestPerEnv(deps []model.Deployment) map[string]model.Deployment {
	out := make(map[string]model.Deployment)
	for _, d := range deps {
		slug := strings.ToLower(strings.TrimSpace(d.EnvironmentSlug))
		if slug == "" {
			slug = "dev"
		}
		if _, seen := out[slug]; !seen {
			out[slug] = d
		}
	}
	return out
}

type Runtime struct {
	store       *store.RuntimeStore
	deployments *Deployments
	catalog     *Catalog
}

func NewRuntime(s *store.RuntimeStore, deployments *Deployments, catalog *Catalog) *Runtime {
	return &Runtime{store: s, deployments: deployments, catalog: catalog}
}

func (r *Runtime) RequestStop(ctx context.Context, serviceID, environment string) (model.RuntimeJob, error) {
	return r.enqueue(ctx, serviceID, environment, "stop", "running")
}

func (r *Runtime) RequestStart(ctx context.Context, serviceID, environment string) (model.RuntimeJob, error) {
	return r.enqueue(ctx, serviceID, environment, "start", "stopped")
}

func (r *Runtime) enqueue(ctx context.Context, serviceID, environment, action, requiredStatus string) (model.RuntimeJob, error) {
	serviceID = strings.TrimSpace(serviceID)
	if serviceID == "" {
		return model.RuntimeJob{}, fmt.Errorf("%w: service_id is required", ErrInvalid)
	}

	slug := strings.ToLower(strings.TrimSpace(environment))
	if slug == "" {
		slug = "dev"
	}

	svc, err := r.catalog.Get(ctx, serviceID)
	if err != nil {
		return model.RuntimeJob{}, err
	}

	deps, err := r.deployments.ListByService(ctx, serviceID)
	if err != nil {
		return model.RuntimeJob{}, err
	}
	if len(deps) == 0 {
		return model.RuntimeJob{}, fmt.Errorf("%w: no deployments for this service", ErrInvalid)
	}

	var target *model.Deployment
	for i := range deps {
		if deps[i].EnvironmentSlug == slug {
			cp := deps[i]
			target = &cp
			break
		}
	}

	if target == nil {
		return model.RuntimeJob{}, fmt.Errorf("%w: no deployment for environment %q", ErrInvalid, slug)
	}

	if target.Status != requiredStatus {
		return model.RuntimeJob{}, fmt.Errorf("%w: %s deployment must be %s (current: %s)", ErrInvalid, slug, requiredStatus, target.Status)
	}

	busy, err := r.store.HasActiveJob(ctx, target.ID)
	if err != nil {
		return model.RuntimeJob{}, err
	}
	if busy {
		return model.RuntimeJob{}, fmt.Errorf("%w: runtime job already in progress for this deployment", ErrInvalid)
	}

	job, err := r.store.Create(ctx, serviceID, target.ID, svc.Name, action)
	if err != nil {
		return model.RuntimeJob{}, fmt.Errorf("enqueue runtime job: %w", err)
	}
	return job, nil
}

func (r *Runtime) ClaimNext(ctx context.Context, workerID string) (model.RuntimeJob, error) {
	workerID = strings.TrimSpace(workerID)
	if workerID == "" {
		return model.RuntimeJob{}, fmt.Errorf("%w: worker_id is required", ErrInvalid)
	}

	job, err := r.store.ClaimNext(ctx, workerID)
	if errors.Is(err, store.ErrNotFound) {
		return model.RuntimeJob{}, ErrNotFound
	}
	if err != nil {
		return model.RuntimeJob{}, fmt.Errorf("claim runtime job: %w", err)
	}
	return job, nil
}

func (r *Runtime) UpdateFromAgent(ctx context.Context, id string, req model.UpdateRuntimeJobRequest) (model.RuntimeJob, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return model.RuntimeJob{}, fmt.Errorf("%w: id is required", ErrInvalid)
	}

	req.Status = strings.TrimSpace(req.Status)
	req.ErrorMessage = strings.TrimSpace(req.ErrorMessage)

	if req.Status != "" {
		switch req.Status {
		case "done", "failed":
		default:
			return model.RuntimeJob{}, fmt.Errorf("%w: invalid status %q", ErrInvalid, req.Status)
		}
	}

	existing, err := r.store.Update(ctx, id, req)
	if err != nil {
		return model.RuntimeJob{}, mapRepoErr(err)
	}

	if req.Status == "done" {
		deployStatus := ""
		switch existing.Action {
		case "stop":
			deployStatus = "stopped"
		case "start":
			deployStatus = "running"
		case "remove":
		}
		if deployStatus != "" {
			_, err := r.deployments.Update(ctx, existing.DeploymentID, model.UpdateDeploymentRequest{
				Status: deployStatus,
			})
			if err != nil && !errors.Is(err, ErrNotFound) {
				return model.RuntimeJob{}, err
			}
		}
	}

	return existing, nil
}

func (r *Runtime) ValidateDelete(ctx context.Context, svc model.Service) error {
	deps, err := r.deployments.ListByService(ctx, svc.ID)
	if err != nil {
		return err
	}
	if len(deps) == 0 {
		return nil
	}

	perEnv := latestPerEnv(deps)
	var running []string
	for slug, d := range perEnv {
		if d.Status == "running" {
			running = append(running, slug)
		}
	}
	if len(running) == 0 {
		return nil
	}
	slices.Sort(running)

	for _, slug := range running {
		if slug == "prod" {
			return fmt.Errorf(
				"%w: production is still running; stop prod before deleting this service",
				ErrForbidden,
			)
		}
	}
	return fmt.Errorf(
		"%w: stop running environments before delete: %s",
		ErrConflict,
		strings.Join(running, ", "),
	)
}

func (r *Runtime) EnqueueRemove(ctx context.Context, svc model.Service) error {
	deps, err := r.deployments.ListByService(ctx, svc.ID)
	if err != nil {
		return err
	}
	if len(deps) == 0 {
		return nil
	}

	perEnv := latestPerEnv(deps)
	for slug, d := range perEnv {
		busy, err := r.store.HasActiveJob(ctx, d.ID)
		if err != nil {
			return err
		}
		if busy {
			return fmt.Errorf(
				"%w: runtime job in progress for environment %q",
				ErrConflict,
				slug,
			)
		}
		if _, err := r.store.Create(ctx, svc.ID, d.ID, svc.Name, "remove"); err != nil {
			return fmt.Errorf("enqueue remove job for %q: %w", slug, err)
		}
	}
	return nil
}
