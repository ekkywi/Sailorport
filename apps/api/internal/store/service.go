package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/ekkywi/sailorport/apps/api/internal/model"
	"github.com/jackc/pgx/v5/pgconn"
)

var ErrNotFound = errors.New("Service not found")
var ErrConflict = errors.New("Service already exists")

type ServicesStore struct {
	db *sql.DB
}

func NewServicesStore(db *sql.DB) *ServicesStore {
	return &ServicesStore{db: db}
}

func nullUUID(id string) any {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil
	}
	return id
}

func scanOwnerUserID(ns sql.NullString) string {
	if !ns.Valid {
		return ""
	}
	return ns.String
}

func (s *ServicesStore) List(ctx context.Context) ([]model.Service, error) {
	const q = `
	SELECT id , name, description, owner, owner_user_id, template_id, workspace_path,
		source_type, repo_url, branch, dockerfile_path,
		webhook_secret, auto_deploy_enabled, auto_deploy_environment,
		catalog_app_id, image, container_port,
		created_at, updated_at
	FROM services
	ORDER BY created_at DESC
	`

	rows, err := s.db.QueryContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("List services: %w", err)
	}
	defer rows.Close()

	services := make([]model.Service, 0)
	for rows.Next() {
		var svc model.Service
		var ownerUserID sql.NullString
		if err := rows.Scan(
			&svc.ID,
			&svc.Name,
			&svc.Description,
			&svc.Owner,
			&ownerUserID,
			&svc.TemplateID,
			&svc.WorkspacePath,
			&svc.SourceType,
			&svc.RepoURL,
			&svc.Branch,
			&svc.DockerfilePath,
			&svc.WebhookSecret,
			&svc.AutoDeployEnabled,
			&svc.AutoDeployEnvironment,
			&svc.CatalogAppID,
			&svc.Image,
			&svc.ContainerPort,
			&svc.CreatedAt,
			&svc.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("Scan service: %w", err)
		}
		svc.OwnerUserID = scanOwnerUserID(ownerUserID)
		services = append(services, svc)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Iterate services: %w", err)
	}
	return services, nil
}

func (s *ServicesStore) ListByOwner(ctx context.Context, ownerUserID string) ([]model.Service, error) {
	ownerUserID = strings.TrimSpace(ownerUserID)
	if ownerUserID == "" {
		return nil, fmt.Errorf("owner_user_id is required")
	}

	const q = `
	SELECT id, name , description, owner, owner_user_id, template_id, workspace_path,
		source_type, repo_url, branch, dockerfile_path,
		webhook_secret, auto_deploy_enabled, auto_deploy_environment,
		catalog_app_id, image, container_port,
		created_at, updated_at
	FROM services
	WHERE owner_user_id = $1
	ORDER BY created_at DESC
	`

	rows, err := s.db.QueryContext(ctx, q, ownerUserID)
	if err != nil {
		return nil, fmt.Errorf("List services by owner: %w", err)
	}
	defer rows.Close()

	services := make([]model.Service, 0)
	for rows.Next() {
		var svc model.Service
		var ownerUserIDCol sql.NullString
		if err := rows.Scan(
			&svc.ID,
			&svc.Name,
			&svc.Description,
			&svc.Owner,
			&ownerUserIDCol,
			&svc.TemplateID,
			&svc.WorkspacePath,
			&svc.SourceType,
			&svc.RepoURL,
			&svc.Branch,
			&svc.DockerfilePath,
			&svc.WebhookSecret,
			&svc.AutoDeployEnabled,
			&svc.AutoDeployEnvironment,
			&svc.CatalogAppID,
			&svc.Image,
			&svc.ContainerPort,
			&svc.CreatedAt,
			&svc.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("Scan service: %w", err)
		}
		svc.OwnerUserID = scanOwnerUserID(ownerUserIDCol)
		services = append(services, svc)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Iterate services: %w", err)
	}
	return services, nil
}

func (s *ServicesStore) Get(ctx context.Context, id string) (model.Service, error) {
	const q = `
	SELECT id, name, description, owner, owner_user_id, template_id, workspace_path,
		source_type, repo_url, branch, dockerfile_path,
		webhook_secret, auto_deploy_enabled, auto_deploy_environment,
		catalog_app_id, image, container_port,
		created_at, updated_at
	FROM services
	WHERE id = $1
	`

	var svc model.Service
	var ownerUserID sql.NullString
	err := s.db.QueryRowContext(ctx, q, id).Scan(
		&svc.ID,
		&svc.Name,
		&svc.Description,
		&svc.Owner,
		&ownerUserID,
		&svc.TemplateID,
		&svc.WorkspacePath,
		&svc.SourceType,
		&svc.RepoURL,
		&svc.Branch,
		&svc.DockerfilePath,
		&svc.WebhookSecret,
		&svc.AutoDeployEnabled,
		&svc.AutoDeployEnvironment,
		&svc.CatalogAppID,
		&svc.Image,
		&svc.ContainerPort,
		&svc.CreatedAt,
		&svc.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return model.Service{}, ErrNotFound
	}
	if err != nil {
		return model.Service{}, fmt.Errorf("Get service: %w", err)
	}
	svc.OwnerUserID = scanOwnerUserID(ownerUserID)
	return svc, nil
}

func (s *ServicesStore) Create(ctx context.Context, req model.CreateServiceRequest) (model.Service, error) {
	const q = `
		INSERT INTO services (
			name, description, owner, owner_user_id, template_id, workspace_path,
			source_type, repo_url, branch, dockerfile_path,
			webhook_secret, auto_deploy_enabled, auto_deploy_environment,
			catalog_app_id, image, container_port
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
		RETURNING id, name, description, owner, owner_user_id, template_id, workspace_path,
			source_type, repo_url, branch, dockerfile_path,
			webhook_secret, auto_deploy_enabled, auto_deploy_environment,
			catalog_app_id, image, container_port,
			created_at, updated_at
	`

	var svc model.Service
	var ownerUserID sql.NullString
	err := s.db.QueryRowContext(ctx, q,
		req.Name,
		req.Description,
		req.Owner,
		nullUUID(req.OwnerUserID),
		req.TemplateID,
		req.WorkspacePath,
		req.SourceType,
		req.RepoURL,
		req.Branch,
		req.DockerfilePath,
		req.WebhookSecret,
		req.AutoDeployEnabled,
		req.AutoDeployEnvironment,
		req.CatalogAppID,
		req.Image,
		req.ContainerPort,
	).Scan(
		&svc.ID,
		&svc.Name,
		&svc.Description,
		&svc.Owner,
		&ownerUserID,
		&svc.TemplateID,
		&svc.WorkspacePath,
		&svc.SourceType,
		&svc.RepoURL,
		&svc.Branch,
		&svc.DockerfilePath,
		&svc.WebhookSecret,
		&svc.AutoDeployEnabled,
		&svc.AutoDeployEnvironment,
		&svc.CatalogAppID,
		&svc.Image,
		&svc.ContainerPort,
		&svc.CreatedAt,
		&svc.UpdatedAt,
	)
	if err != nil {
		if isUniqueViolation(err) {
			return model.Service{}, ErrConflict
		}
		return model.Service{}, fmt.Errorf("Create service: %w", err)
	}
	svc.OwnerUserID = scanOwnerUserID(ownerUserID)
	return svc, nil
}

func (s *ServicesStore) Update(ctx context.Context, id string, req model.UpdateServiceRequest) (model.Service, error) {
	const q = `
		UPDATE services
		SET name = $1,
			description = $2,
			owner = $3,
			source_type = $4,
			repo_url = $5,
			branch = $6,
			dockerfile_path = $7,
			webhook_secret = $8,
			auto_deploy_enabled = $9,
			auto_deploy_environment = $10,
			catalog_app_id = $11,
			image = $12,
			container_port = $13,
			updated_at = NOW()
		WHERE id = $14
		RETURNING id, name, description, owner, owner_user_id, template_id, workspace_path,
			source_type, repo_url, branch, dockerfile_path,
			webhook_secret, auto_deploy_enabled, auto_deploy_environment,
			catalog_app_id, image, container_port,
			created_at, updated_at
	`

	var svc model.Service
	var ownerUserID sql.NullString
	err := s.db.QueryRowContext(ctx, q,
		req.Name,
		req.Description,
		req.Owner,
		req.SourceType,
		req.RepoURL,
		req.Branch,
		req.DockerfilePath,
		req.WebhookSecret,
		*req.AutoDeployEnabled,
		req.AutoDeployEnvironment,
		req.CatalogAppID,
		req.Image,
		req.ContainerPort,
		id,
	).Scan(
		&svc.ID,
		&svc.Name,
		&svc.Description,
		&svc.Owner,
		&ownerUserID,
		&svc.TemplateID,
		&svc.WorkspacePath,
		&svc.SourceType,
		&svc.RepoURL,
		&svc.Branch,
		&svc.DockerfilePath,
		&svc.WebhookSecret,
		&svc.AutoDeployEnabled,
		&svc.AutoDeployEnvironment,
		&svc.CatalogAppID,
		&svc.Image,
		&svc.ContainerPort,
		&svc.CreatedAt,
		&svc.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return model.Service{}, ErrNotFound
	}
	if err != nil {
		if isUniqueViolation(err) {
			return model.Service{}, ErrConflict
		}
		return model.Service{}, fmt.Errorf("update service: %w", err)
	}
	svc.OwnerUserID = scanOwnerUserID(ownerUserID)
	return svc, nil
}

func (s *ServicesStore) Delete(ctx context.Context, id string) error {
	const q = `DELETE FROM services WHERE id = $1`

	result, err := s.db.ExecContext(ctx, q, id)
	if err != nil {
		return fmt.Errorf("Delete service: %w", err)
	}
	n, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Rows Affected: %w", err)
	}
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}
	return strings.Contains(err.Error(), "SQLSTATE 23505")
}

func (s *ServicesStore) UpdateOwner(ctx context.Context, id, ownerUserID, ownerLabel string) (model.Service, error) {
	id = strings.TrimSpace(id)
	ownerUserID = strings.TrimSpace(ownerUserID)
	ownerLabel = strings.TrimSpace(ownerLabel)
	if id == "" || ownerUserID == "" {
		return model.Service{}, fmt.Errorf("id and owner_user_id are required")
	}
	if ownerLabel == "" {
		ownerLabel = ownerUserID
	}

	const q = `
		UPDATE services
		SET owner = $1,
			owner_user_id = $2,
			updated_at = NOW()
		WHERE id = $3
		RETURNING id, name, description, owner, owner_user_id, template_id, workspace_path,
			source_type, repo_url, branch, dockerfile_path,
			webhook_secret, auto_deploy_enabled, auto_deploy_environment,
			catalog_app_id, image, container_port,
			created_at, updated_at
	`

	var svc model.Service
	var ownerUserIDCol sql.NullString
	err := s.db.QueryRowContext(ctx, q, ownerLabel, ownerUserID, id).Scan(
		&svc.ID,
		&svc.Name,
		&svc.Description,
		&svc.Owner,
		&ownerUserIDCol,
		&svc.TemplateID,
		&svc.WorkspacePath,
		&svc.SourceType,
		&svc.RepoURL,
		&svc.Branch,
		&svc.DockerfilePath,
		&svc.WebhookSecret,
		&svc.AutoDeployEnabled,
		&svc.AutoDeployEnvironment,
		&svc.CatalogAppID,
		&svc.Image,
		&svc.ContainerPort,
		&svc.CreatedAt,
		&svc.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return model.Service{}, ErrNotFound
	}
	if err != nil {
		return model.Service{}, fmt.Errorf("Update owner: %w", err)
	}
	svc.OwnerUserID = scanOwnerUserID(ownerUserIDCol)
	return svc, nil
}
