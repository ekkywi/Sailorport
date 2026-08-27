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
	envs    *store.EnvironmentsStore
	catalog *Catalog
	workers *Workers
}

func NewDeployments(s *store.DeploymentsStore, envs *store.EnvironmentsStore, catalog *Catalog, workers *Workers) *Deployments {
	return &Deployments{store: s, envs: envs, catalog: catalog, workers: workers}
}

func (d *Deployments) Create(ctx context.Context, serviceID string, req model.CreateDeploymentRequest) (model.Deployment, error) {
	serviceID = strings.TrimSpace(serviceID)
	if serviceID == "" {
		return model.Deployment{}, fmt.Errorf("%w: service_id is required", ErrInvalid)
	}

	svc, err := d.catalog.Get(ctx, serviceID)
	if err != nil {
		return model.Deployment{}, err
	}
	hasWorkspace := strings.TrimSpace(svc.WorkspacePath) != ""
	isGit := svc.SourceType == "git" && strings.TrimSpace(svc.RepoURL) != ""
	isCatalogApp := svc.SourceType == "catalog_app" && strings.TrimSpace(svc.Image) != ""
	if !hasWorkspace && !isGit && !isCatalogApp {
		return model.Deployment{}, fmt.Errorf(
			"%w: service needs a workspace (scaffold) or git repo_url",
			ErrInvalid,
		)
	}

	slug := strings.ToLower(strings.TrimSpace(req.Environment))
	if slug == "" {
		slug = "dev"
	}

	env, err := d.envs.GetBySlug(ctx, slug)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return model.Deployment{}, fmt.Errorf("%w: unknown environment %q", ErrInvalid, slug)
		}
		return model.Deployment{}, err
	}

	targetWorkerID, err := d.resolveTargetWorker(ctx, serviceID, slug, req.WorkerID)
	if err != nil {
		return model.Deployment{}, err
	}

	gitSHA := strings.TrimSpace(req.GitSHA)
	if gitSHA != "" && svc.SourceType != "git" {
		return model.Deployment{}, fmt.Errorf(
			"%w: git_sha is only valid for git services",
			ErrInvalid,
		)
	}

	out, err := d.store.Create(ctx, serviceID, env.ID, targetWorkerID, gitSHA)
	if err != nil {
		return model.Deployment{}, fmt.Errorf("Create deployment: %w", err)
	}
	return out, nil
}

func (d *Deployments) Redeploy(ctx context.Context, deploymentID string) (model.Deployment, error) {
	deploymentID = strings.TrimSpace(deploymentID)
	if deploymentID == "" {
		return model.Deployment{}, fmt.Errorf("%w: id is required", ErrInvalid)
	}

	src, err := d.store.Get(ctx, deploymentID)
	if err != nil {
		return model.Deployment{}, mapRepoErr(err)
	}

	sha := strings.TrimSpace(src.GitSHA)
	if sha == "" {
		return model.Deployment{}, fmt.Errorf(
			"%w: deployment has no git_sha to redeploy",
			ErrInvalid,
		)
	}

	return d.Create(ctx, src.ServiceID, model.CreateDeploymentRequest{
		Environment: src.EnvironmentSlug,
		GitSHA:      sha,
	})
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
	req.GitSHA = strings.TrimSpace(req.GitSHA)
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

func (d *Deployments) resolveTargetWorker(ctx context.Context, serviceID, envSlug, requestedWorkerID string) (*string, error) {
	requestedWorkerID = strings.TrimSpace(requestedWorkerID)

	if requestedWorkerID != "" {
		return d.validateWorkerForDeploy(ctx, requestedWorkerID, envSlug, true)
	}

	deps, err := d.store.ListByService(ctx, serviceID)
	if err != nil {
		return nil, err
	}
	for _, dep := range deps {
		if dep.EnvironmentSlug != envSlug {
			continue
		}
		if dep.WorkerID != nil && strings.TrimSpace(*dep.WorkerID) != "" {
			id := strings.TrimSpace(*dep.WorkerID)
			target, err := d.validateWorkerForDeploy(ctx, id, envSlug, false)
			if errors.Is(err, ErrNotFound) {
				return nil, nil
			}
			return target, err
		}
		break
	}
	return nil, nil
}

func (d *Deployments) validateWorkerForDeploy(
	ctx context.Context,
	workerID, envSlug string,
	requireOnline bool,
) (*string, error) {
	w, err := d.workers.Get(ctx, workerID)
	if err != nil {
		return nil, err
	}
	if requireOnline && w.Status != "online" {
		return nil, fmt.Errorf(
			"%w: worker %q is %s (must be online)",
			ErrConflict, w.Name, w.Status,
		)
	}
	if err := workerEnvConflict(w, envSlug); err != nil {
		return nil, err
	}
	id := w.ID
	return &id, nil
}
