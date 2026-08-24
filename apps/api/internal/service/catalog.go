package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/ekkywi/sailorport/apps/api/internal/model"
	"github.com/ekkywi/sailorport/apps/api/internal/store"
)

type CleanupEnqueue interface {
	ValidateDelete(ctx context.Context, svc model.Service) error
	EnqueueRemove(ctx context.Context, svc model.Service) error
}

var (
	ErrInvalid      = errors.New("invalid input")
	ErrNotFound     = errors.New("service not found")
	ErrConflict     = errors.New("service already exists")
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")
)

type Repository interface {
	List(ctx context.Context) ([]model.Service, error)
	Get(ctx context.Context, id string) (model.Service, error)
	Create(ctx context.Context, req model.CreateServiceRequest) (model.Service, error)
	Update(ctx context.Context, id string, req model.UpdateServiceRequest) (model.Service, error)
	Delete(ctx context.Context, id string) error
}

type DeploymentReader interface {
	LatestByServices(ctx context.Context) (map[string]model.Deployment, error)
	LatestPerEnvByServices(ctx context.Context) (map[string]map[string]model.Deployment, error)
}

type Catalog struct {
	repo         Repository
	deployments  DeploymentReader
	workspaceDir string
	cleanup      CleanupEnqueue
	audit        *Audit
}

func (c *Catalog) SetCleanupEnqueue(e CleanupEnqueue) {
	c.cleanup = e
}

func (c *Catalog) SetAudit(a *Audit) {
	c.audit = a
}

func NewCatalog(repo Repository, deployments DeploymentReader, workspaceDir string) *Catalog {
	return &Catalog{
		repo:         repo,
		deployments:  deployments,
		workspaceDir: workspaceDir,
	}
}

func (c *Catalog) List(ctx context.Context) ([]model.Service, error) {
	services, err := c.repo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("list services: %w", err)
	}
	if c.deployments == nil || len(services) == 0 {
		return services, nil
	}
	latest, err := c.deployments.LatestByServices(ctx)
	if err != nil {
		return nil, fmt.Errorf("List services latest deployments: %w", err)
	}
	perEnv, err := c.deployments.LatestPerEnvByServices(ctx)
	if err != nil {
		return nil, fmt.Errorf("List services env deployments: %w", err)
	}
	for i := range services {
		if d, ok := latest[services[i].ID]; ok {
			cp := d
			services[i].LatestDeployment = &cp
		}
		if envs, ok := perEnv[services[i].ID]; ok && len(envs) > 0 {
			services[i].EnvDeployments = envs
		}
	}
	return services, nil
}

func (c *Catalog) Get(ctx context.Context, id string) (model.Service, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return model.Service{}, fmt.Errorf("%w: id is required", ErrInvalid)
	}
	svc, err := c.repo.Get(ctx, id)
	if err != nil {
		return model.Service{}, mapRepoErr(err)
	}
	return svc, nil
}

func (c *Catalog) Create(ctx context.Context, req model.CreateServiceRequest, actorID, actorEmail string) (model.Service, error) {
	req, err := normalizeCreate(req)
	if err != nil {
		return model.Service{}, err
	}
	svc, err := c.repo.Create(ctx, req)
	if err != nil {
		return model.Service{}, mapRepoErr(err)
	}
	c.recordService(ctx, actorID, actorEmail, "service.create", svc)
	return svc, nil
}

func (c *Catalog) Update(ctx context.Context, id string, req model.UpdateServiceRequest, actorID, actorEmail string) (model.Service, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return model.Service{}, fmt.Errorf("%w: id is required", ErrInvalid)
	}

	existing, err := c.repo.Get(ctx, id)
	if err != nil {
		return model.Service{}, mapRepoErr(err)
	}

	req, err = normalizeUpdate(req, existing)
	if err != nil {
		return model.Service{}, err
	}
	svc, err := c.repo.Update(ctx, id, req)
	if err != nil {
		return model.Service{}, mapRepoErr(err)
	}
	c.recordService(ctx, actorID, actorEmail, "service.update", svc)
	return svc, nil
}

func (c *Catalog) Delete(ctx context.Context, id, actorID, actorEmail string) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return fmt.Errorf("%w: id is required", ErrInvalid)
	}

	svc, err := c.repo.Get(ctx, id)
	if err != nil {
		return mapRepoErr(err)
	}

	if c.cleanup != nil {
		if err := c.cleanup.ValidateDelete(ctx, svc); err != nil {
			return err
		}
		if err := c.cleanup.EnqueueRemove(ctx, svc); err != nil {
			return fmt.Errorf("enqueue container cleanup: %w", err)
		}
	}

	if err := c.recordDelete(ctx, svc, actorID, actorEmail); err != nil {
		return err
	}

	if err := c.repo.Delete(ctx, id); err != nil {
		return mapRepoErr(err)
	}

	if err := c.removeWorkspace(svc.WorkspacePath); err != nil {
		log.Printf("catalog delete: workspace cleanup failed for %s: %v", svc.Name, err)
		return fmt.Errorf("deleted from catalog but workspace cleanup failed: %w", err)
	}
	return nil
}

func (c *Catalog) recordDelete(ctx context.Context, svc model.Service, actorID, actorEmail string) error {
	if c.audit == nil {
		return nil
	}
	payload := map[string]any{
		"description":    svc.Description,
		"owner":          svc.Owner,
		"template_id":    svc.TemplateID,
		"workspace_path": svc.WorkspacePath,
	}
	if c.deployments != nil {
		perEnv, err := c.deployments.LatestPerEnvByServices(ctx)
		if err == nil {
			if m, ok := perEnv[svc.ID]; ok && len(m) > 0 {
				slugs := make([]string, 0, len(m))
				statuses := make(map[string]string, len(m))
				for slug, d := range m {
					slugs = append(slugs, slug)
					statuses[slug] = d.Status
				}
				slices.Sort(slugs)
				payload["environments"] = slugs
				payload["env_status"] = statuses
			}
		}
	}
	return c.audit.Record(ctx, model.AuditRecord{
		ActorID:      strings.TrimSpace(actorID),
		ActorEmail:   strings.TrimSpace(actorEmail),
		Action:       "service.delete",
		ResourceType: "service",
		ResourceID:   svc.ID,
		ResourceName: svc.Name,
		Payload:      payload,
	})
}

func (c *Catalog) recordService(ctx context.Context, actorID, actorEmail, action string, svc model.Service) {
	if c.audit == nil {
		return
	}
	err := c.audit.Record(ctx, model.AuditRecord{
		ActorID:      strings.TrimSpace(actorID),
		ActorEmail:   strings.TrimSpace(actorEmail),
		Action:       action,
		ResourceType: "service",
		ResourceID:   svc.ID,
		ResourceName: svc.Name,
		Payload: map[string]any{
			"description":     svc.Description,
			"owner":           svc.Owner,
			"template_id":     svc.TemplateID,
			"workspace_path":  svc.WorkspacePath,
			"source_type":     svc.SourceType,
			"repo_url":        svc.RepoURL,
			"branch":          svc.Branch,
			"dockerfile_path": svc.DockerfilePath,
		},
	})
	if err != nil {
		log.Printf("audit %s: %v", action, err)
	}
}

func (c *Catalog) removeWorkspace(workspacePath string) error {
	workspacePath = strings.TrimSpace(workspacePath)
	if workspacePath == "" || strings.TrimSpace(c.workspaceDir) == "" {
		return nil
	}

	cleanRoot, err := filepath.Abs(filepath.Clean(c.workspaceDir))
	if err != nil {
		return err
	}
	cleanPath, err := filepath.Abs(filepath.Clean(workspacePath))
	if err != nil {
		return err
	}

	sep := string(filepath.Separator)
	if cleanPath == cleanRoot || !strings.HasPrefix(cleanPath+sep, cleanRoot+sep) {
		return nil
	}

	if err := os.RemoveAll(cleanPath); err != nil {
		return err
	}
	return nil
}

func normalizeCreate(req model.CreateServiceRequest) (model.CreateServiceRequest, error) {
	req.Name = strings.TrimSpace(req.Name)
	req.Description = strings.TrimSpace(req.Description)
	req.Owner = strings.TrimSpace(req.Owner)
	req.TemplateID = strings.TrimSpace(req.TemplateID)
	req.WorkspacePath = strings.TrimSpace(req.WorkspacePath)
	req.SourceType = strings.TrimSpace(req.SourceType)
	req.RepoURL = strings.TrimSpace(req.RepoURL)
	req.Branch = strings.TrimSpace(req.Branch)
	req.DockerfilePath = strings.TrimSpace(req.DockerfilePath)
	req.WebhookSecret = strings.TrimSpace(req.WebhookSecret)
	req.AutoDeployEnvironment = strings.TrimSpace(req.AutoDeployEnvironment)

	if req.SourceType == "" {
		req.SourceType = "scaffold"
	}
	if req.Branch == "" {
		req.Branch = "main"
	}
	if req.DockerfilePath == "" {
		req.DockerfilePath = "Dockerfile"
	}
	if req.AutoDeployEnvironment == "" {
		req.AutoDeployEnvironment = "staging"
	}

	if req.Name == "" {
		return req, fmt.Errorf("%w: name is required", ErrInvalid)
	}
	if err := validateSourceFields(req.SourceType, req.RepoURL); err != nil {
		return req, err
	}
	if err := validateAutoDeployEnvironment(req.AutoDeployEnvironment); err != nil {
		return req, err
	}
	return req, nil
}

func normalizeUpdate(req model.UpdateServiceRequest, existing model.Service) (model.UpdateServiceRequest, error) {
	req.Name = strings.TrimSpace(req.Name)
	req.Description = strings.TrimSpace(req.Description)
	req.Owner = strings.TrimSpace(req.Owner)
	req.SourceType = strings.TrimSpace(req.SourceType)
	req.RepoURL = strings.TrimSpace(req.RepoURL)
	req.Branch = strings.TrimSpace(req.Branch)
	req.DockerfilePath = strings.TrimSpace(req.DockerfilePath)
	req.WebhookSecret = strings.TrimSpace(req.WebhookSecret)
	req.AutoDeployEnvironment = strings.TrimSpace(req.AutoDeployEnvironment)

	if req.Name == "" {
		return req, fmt.Errorf("%w: name is required", ErrInvalid)
	}

	if req.SourceType == "" {
		req.SourceType = existing.SourceType
	}
	if req.RepoURL == "" {
		req.RepoURL = existing.RepoURL
	}
	if req.Branch == "" {
		req.Branch = existing.Branch
	}
	if req.DockerfilePath == "" {
		req.DockerfilePath = existing.DockerfilePath
	}
	if req.Branch == "" {
		req.Branch = "main"
	}
	if req.DockerfilePath == "" {
		req.DockerfilePath = "Dockerfile"
	}
	if req.WebhookSecret == "" {
		req.WebhookSecret = existing.WebhookSecret
	}
	if req.AutoDeployEnvironment == "" {
		req.AutoDeployEnvironment = existing.AutoDeployEnvironment
	}
	if req.AutoDeployEnvironment == "" {
		req.AutoDeployEnvironment = "staging"
	}
	if req.AutoDeployEnabled == nil {
		v := existing.AutoDeployEnabled
		req.AutoDeployEnabled = &v
	}

	if err := validateSourceFields(req.SourceType, req.RepoURL); err != nil {
		return req, err
	}
	if err := validateAutoDeployEnvironment(req.AutoDeployEnvironment); err != nil {
		return req, err
	}
	return req, nil
}

func validateSourceFields(sourceType, repoURL string) error {
	switch sourceType {
	case "scaffold", "git":
	default:
		return fmt.Errorf("%w: invalid source_type %q", ErrInvalid, sourceType)
	}
	if sourceType == "git" && repoURL == "" {
		return fmt.Errorf("%w: repo_url is required for source_type=git", ErrInvalid)
	}
	return nil
}

func mapRepoErr(err error) error {
	if errors.Is(err, store.ErrNotFound) {
		return ErrNotFound
	}
	if errors.Is(err, store.ErrConflict) {
		return ErrConflict
	}
	return err
}

func validateAutoDeployEnvironment(env string) error {
	switch env {
	case "dev", "staging", "prod":
		return nil
	default:
		return fmt.Errorf("%w: invalid auto_deploy_environment %q (want dev, staging, or prod)", ErrInvalid, env)
	}
}
