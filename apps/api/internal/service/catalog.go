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

	"github.com/ekkywi/sailorport/apps/api/internal/catalogapp"
	"github.com/ekkywi/sailorport/apps/api/internal/model"
	"github.com/ekkywi/sailorport/apps/api/internal/secrets"
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
	apps         *catalogapp.Registry
	secrets      secrets.Store
}

func (c *Catalog) SetCleanupEnqueue(e CleanupEnqueue) {
	c.cleanup = e
}

func (c *Catalog) SetAudit(a *Audit) {
	c.audit = a
}

func NewCatalog(
	repo Repository,
	deployments DeploymentReader,
	workspaceDir string,
	apps *catalogapp.Registry,
	secretsStore secrets.Store,
) *Catalog {
	return &Catalog{
		repo:         repo,
		deployments:  deployments,
		workspaceDir: workspaceDir,
		apps:         apps,
		secrets:      secretsStore,
	}
}

func (c *Catalog) List(ctx context.Context) ([]model.Service, error) {
	services, err := c.repo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("list services: %w", err)
	}
	if len(services) == 0 {
		return services, nil
	}
	if c.deployments != nil {
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
	}
	for i := range services {
		if err := c.attachCatalogEnvPublic(ctx, &services[i]); err != nil {
			return nil, err
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
	if err := c.attachCatalogEnvPublic(ctx, &svc); err != nil {
		return model.Service{}, err
	}
	return svc, nil
}

func (c *Catalog) Create(ctx context.Context, req model.CreateServiceRequest, actorID, actorEmail string) (model.Service, error) {
	req, err := normalizeCreate(req)
	if err != nil {
		return model.Service{}, err
	}
	req, err = c.applyCatalogAppDefaults(req)
	if err != nil {
		return model.Service{}, err
	}

	var envEntries []model.CatalogEnv
	if req.SourceType == "catalog_app" {
		m, err := c.catalogAppManifest(req.CatalogAppID)
		if err != nil {
			return model.Service{}, err
		}
		envEntries, err = buildCatalogEnvEntries(m, req.CatalogEnv)
		if err != nil {
			return model.Service{}, err
		}
	} else if len(req.CatalogEnv) > 0 {
		return model.Service{}, fmt.Errorf("%w: catalog_env is only allowed for source_type=catalog_app", ErrInvalid)
	}

	svc, err := c.repo.Create(ctx, req)
	if err != nil {
		return model.Service{}, mapRepoErr(err)
	}

	if len(envEntries) > 0 {
		if c.secrets == nil {
			_ = c.repo.Delete(ctx, svc.ID)
			return model.Service{}, fmt.Errorf("catalog secrets store not configured")
		}
		if err := c.secrets.ReplaceAll(ctx, svc.ID, envEntries); err != nil {
			_ = c.repo.Delete(ctx, svc.ID)
			return model.Service{}, fmt.Errorf("save catalog env: %w", err)
		}
	}

	if err := c.attachCatalogEnvPublic(ctx, &svc); err != nil {
		return model.Service{}, err
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
	if req.CatalogEnv != nil {
		if existing.SourceType != "catalog_app" {
			return model.Service{}, fmt.Errorf("%w: catalog_env is only allowed for source_type=catalog_app", ErrInvalid)
		}
		m, err := c.catalogAppManifest(existing.CatalogAppID)
		if err != nil {
			return model.Service{}, err
		}
		if c.secrets == nil {
			return model.Service{}, fmt.Errorf("catalog secrets store not configured")
		}

		existingPlain, err := c.secrets.ResolveForDeploy(ctx, id)
		if err != nil {
			return model.Service{}, fmt.Errorf("load catalog env: %w", err)
		}
		existingRows := catalogEnvRowsFromMap(existingPlain, m)

		envEntries, err := mergeCatalogEnvForUpdate(m, req.CatalogEnv, existingRows)
		if err != nil {
			return model.Service{}, err
		}
		if err := c.secrets.ReplaceAll(ctx, id, envEntries); err != nil {
			return model.Service{}, fmt.Errorf("save catalog env: %w", err)
		}
	}
	svc, err := c.repo.Update(ctx, id, req)
	if err != nil {
		return model.Service{}, mapRepoErr(err)
	}

	if err := c.attachCatalogEnvPublic(ctx, &svc); err != nil {
		return model.Service{}, err
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
	req.CatalogAppID = strings.TrimSpace(req.CatalogAppID)
	req.Image = strings.TrimSpace(req.Image)

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
	if err := validateSourceFields(req.SourceType, req.RepoURL, req.CatalogAppID); err != nil {
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
	req.CatalogAppID = strings.TrimSpace(req.CatalogAppID)
	req.Image = strings.TrimSpace(req.Image)

	if req.Name == "" {
		return req, fmt.Errorf("%w: name is required", ErrInvalid)
	}

	if req.SourceType == "" {
		req.SourceType = existing.SourceType
	}
	if req.SourceType == "catalog_app" && existing.SourceType != "catalog_app" {
		return req, fmt.Errorf("%w: create a new service for catalog_app (udpate source_type no supported yet)", ErrInvalid)
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
	if req.CatalogAppID == "" {
		req.CatalogAppID = existing.CatalogAppID
	}
	if req.Image == "" {
		req.Image = existing.Image
	}
	if req.ContainerPort == 0 {
		req.ContainerPort = existing.ContainerPort
	}

	if err := validateSourceFields(req.SourceType, req.RepoURL, req.CatalogAppID); err != nil {
		return req, err
	}
	if err := validateAutoDeployEnvironment(req.AutoDeployEnvironment); err != nil {
		return req, err
	}
	return req, nil
}

func validateSourceFields(sourceType, repoURL, catalogAppID string) error {
	switch sourceType {
	case "scaffold", "git", "catalog_app":
	default:
		return fmt.Errorf("%w: invalid source_type %q", ErrInvalid, sourceType)
	}
	if sourceType == "git" && repoURL == "" {
		return fmt.Errorf("%w: repo_url is required for source_type=git", ErrInvalid)
	}
	if sourceType == "catalog_app" && catalogAppID == "" {
		return fmt.Errorf("%w: catalog_app_id is required for source_type=catalog_app", ErrInvalid)
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

func (c *Catalog) applyCatalogAppDefaults(req model.CreateServiceRequest) (model.CreateServiceRequest, error) {
	if req.SourceType != "catalog_app" {
		return req, nil
	}
	if c.apps == nil {
		return req, fmt.Errorf("%w: catalog apps registry not configured", ErrInvalid)
	}
	if req.CatalogAppID == "" {
		return req, fmt.Errorf("%w: catalog_app_id is required for source_type=catalog_app", ErrInvalid)
	}

	m, err := c.apps.Get(req.CatalogAppID)
	if err != nil {
		return req, fmt.Errorf("%w: unknown catalog_app_id %q", ErrInvalid, req.CatalogAppID)
	}

	req.Image, err = resolveCatalogAppImage(m, req.Image)
	if err != nil {
		return req, err
	}
	req.ContainerPort = m.ContainerPort
	req.CatalogAppID = m.ID
	req.RepoURL = ""
	req.Branch = ""
	req.DockerfilePath = ""
	req.TemplateID = ""
	req.WebhookSecret = ""
	req.AutoDeployEnabled = false

	return req, nil
}

func (c *Catalog) catalogAppManifest(catalogAppID string) (catalogapp.Manifest, error) {
	if c.apps == nil {
		return catalogapp.Manifest{}, fmt.Errorf("%w: catalog apps registry not configured", ErrInvalid)
	}
	m, err := c.apps.Get(catalogAppID)
	if err != nil {
		return catalogapp.Manifest{}, fmt.Errorf("%w: unknown catalog_app_id %q", ErrInvalid, catalogAppID)
	}
	return m, nil
}

func resolveCatalogAppImage(m catalogapp.Manifest, clientImage string) (string, error) {
	clientImage = strings.TrimSpace(clientImage)
	if len(m.Versions) == 0 {
		return m.Image, nil
	}
	if clientImage == "" {
		return m.Image, nil
	}
	for _, v := range m.Versions {
		if clientImage == v.Image {
			return clientImage, nil
		}
	}
	return "", fmt.Errorf("%w: image %q is not allowed for catalog_app_id %q", ErrInvalid, clientImage, m.ID)
}

func (c *Catalog) listCatalogEnvPlaintext(ctx context.Context, serviceID string) ([]model.CatalogEnv, error) {
	if c.secrets == nil {
		return nil, fmt.Errorf("catalog secrets store not configured")
	}
	rows, err := c.secrets.ResolveForDeploy(ctx, serviceID)
	if err != nil {
		return nil, err
	}
	out := make([]model.CatalogEnv, 0, len(rows))
	for k, v := range rows {
		out = append(out, model.CatalogEnv{Key: k, Value: v})
	}
	return out, nil
}

func catalogEnvRowsFromMap(m map[string]string, manifest catalogapp.Manifest) []model.CatalogEnv {
	secretByName := map[string]bool{}
	for _, f := range manifest.Env {
		secretByName[f.Name] = f.Secret
	}
	out := make([]model.CatalogEnv, 0, len(m))
	for k, v := range m {
		out = append(out, model.CatalogEnv{
			Key: k, Value: v, Secret: secretByName[k],
		})
	}
	return out
}
