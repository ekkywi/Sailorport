package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/ekkywi/sailorport/apps/api/internal/model"
	"github.com/ekkywi/sailorport/apps/api/internal/store"
)

type Runtime struct {
	store       *store.RuntimeStore
	deployments *Deployments
	catalog     *Catalog
}

func NewRuntime(s *store.RuntimeStore, deployments *Deployments, catalog *Catalog) *Runtime {
	return &Runtime{store: s, deployments: deployments, catalog: catalog}
}

func (r *Runtime) RequestStop(ctx context.Context, serviceID string) (model.RuntimeJob, error) {
	return r.enqueue(ctx, serviceID, "stop", "running")
}

func (r *Runtime) RequestStart(ctx context.Context, serviceID string) (model.RuntimeJob, error) {
	return r.enqueue(ctx, serviceID, "start", "stopped")
}

func (r *Runtime) enqueue(ctx context.Context, serviceID, action, requiredStatus string) (model.RuntimeJob, error) {
	serviceID = strings.TrimSpace(serviceID)
	if serviceID == "" {
		return model.RuntimeJob{}, fmt.Errorf("%w: service_id is required", ErrInvalid)
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

	latest := deps[0]
	if latest.Status != requiredStatus {
		return model.RuntimeJob{}, fmt.Errorf("%w: deployment must be %s (current: %s)", ErrInvalid, requiredStatus, latest.Status)
	}

	job, err := r.store.Create(ctx, serviceID, latest.ID, svc.Name, action)
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

func (r *Runtime) EnqueueRemove(ctx context.Context, svc model.Service) error {
	deps, err := r.deployments.ListByService(ctx, svc.ID)
	if err != nil {
		return err
	}
	if len(deps) == 0 {
		return nil
	}

	latest := deps[0]
	_, err = r.store.Create(ctx, svc.ID, latest.ID, svc.Name, "remove")
	if err != nil {
		return fmt.Errorf("Enqueue remove job: %w", err)
	}
	return nil
}
