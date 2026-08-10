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

type Deployments struct {
	store   *store.DeploymentsStore
	catalog *Catalog
}

func NewDeployments(s *store.DeploymentsStore, catalog *Catalog) *Deployments {
	return &Deployments{store: s, catalog: catalog}
}

func (d *Deployments) Create(ctx context.Context, serviceID string) (model.Deployment, error) {
	serviceID = strings.TrimSpace(serviceID)
	if serviceID == "" {
		return model.Deployment{}, fmt.Errorf("%w: service_id is required", ErrInvalid)
	}

	svc, err := d.catalog.Get(ctx, serviceID)
	if err != nil {
		return model.Deployment{}, err
	}
	if strings.TrimSpace(svc.WorkspacePath) == "" {
		return model.Deployment{}, fmt.Errorf("%w: service has no workspace (scaffold first)", ErrInvalid)
	}

	out, err := d.store.Create(ctx, serviceID)
	if err != nil {
		return model.Deployment{}, fmt.Errorf("Create deployment: %w", err)
	}
	return out, nil
}

func (d *Deployments) Get(ctx context.Context, id string) (model.Deployment, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return model.Deployment{}, fmt.Errorf("%w: id is required", ErrInvalid)
	}
	out, err := d.store.Get(ctx, id)
	if err != nil {
		return model.Deployment{}, mapRepoErr(err)
	}
	return out, nil
}

func (d *Deployments) List(ctx context.Context) ([]model.Deployment, error) {
	return d.store.List(ctx)
}

func (d *Deployments) ListByService(ctx context.Context, serviceID string) ([]model.Deployment, error) {
	serviceID = strings.TrimSpace(serviceID)
	if serviceID == "" {
		return nil, fmt.Errorf("%w: service_id is required", ErrInvalid)
	}
	if _, err := d.catalog.Get(ctx, serviceID); err != nil {
		return nil, err
	}
	return d.store.ListByService(ctx, serviceID)
}

func (d *Deployments) ClaimNext(ctx context.Context, workerID string) (model.DeploymentJob, error) {
	workerID = strings.TrimSpace(workerID)
	if workerID == "" {
		return model.DeploymentJob{}, fmt.Errorf("%w: worker_id is required", ErrInvalid)
	}

	job, err := d.store.ClaimNext(ctx, workerID)
	if errors.Is(err, sql.ErrNoRows) {
		return model.DeploymentJob{}, ErrNotFound
	}
	if err != nil {
		return model.DeploymentJob{}, fmt.Errorf("Claim job: %w", err)
	}
	return job, nil
}

func (d *Deployments) Update(ctx context.Context, id string, req model.UpdateDeploymentRequest) (model.Deployment, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return model.Deployment{}, fmt.Errorf("%w: id is required", ErrInvalid)
	}

	req.Status = strings.TrimSpace(req.Status)
	req.ImageTag = strings.TrimSpace(req.ImageTag)
	req.ContainerID = strings.TrimSpace(req.ContainerID)
	req.ErrorMessage = strings.TrimSpace(req.ErrorMessage)
	req.WorkerID = strings.TrimSpace(req.WorkerID)

	if req.Status != "" {
		switch req.Status {
		case "claimed", "building", "running", "failed", "stopped":
		default:
			return model.Deployment{}, fmt.Errorf("%w: invalid status %q", ErrInvalid, req.Status)
		}
	}

	out, err := d.store.Update(ctx, id, req)
	if err != nil {
		return model.Deployment{}, mapRepoErr(err)
	}
	return out, nil
}
