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

func (s *ServicesStore) List(ctx context.Context) ([]model.Service, error) {
	const q = `
	SELECT id, name, description, owner, template_id, workspace_path,
       source_type, repo_url, branch, dockerfile_path,
       webhook_secret, auto_deploy_enabled, auto_deploy_environment,
       created_at, updated_at
	FROM services
	ORDER BY created_at DESC`

	rows, err := s.db.QueryContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("List services: %w", err)
	}
	defer rows.Close()

	services := make([]model.Service, 0)
	for rows.Next() {
		var svc model.Service
		if err := rows.Scan(
			&svc.ID,
			&svc.Name,
			&svc.Description,
			&svc.Owner,
			&svc.TemplateID,
			&svc.WorkspacePath,
			&svc.SourceType,
			&svc.RepoURL,
			&svc.Branch,
			&svc.DockerfilePath,
			&svc.WebhookSecret,
			&svc.AutoDeployEnabled,
			&svc.AutoDeployEnvironment,
			&svc.CreatedAt,
			&svc.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("Scan service: %w", err)
		}
		services = append(services, svc)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Iterate services: %w", err)
	}
	return services, nil
}

func (s *ServicesStore) Get(ctx context.Context, id string) (model.Service, error) {
	const q = `
		SELECT id, name, description, owner, template_id, workspace_path,
		       source_type, repo_url, branch, dockerfile_path,
		       webhook_secret, auto_deploy_enabled, auto_deploy_environment,
		       created_at, updated_at
		FROM services
		WHERE id = $1`

	var svc model.Service
	err := s.db.QueryRowContext(ctx, q, id).Scan(
		&svc.ID,
		&svc.Name,
		&svc.Description,
		&svc.Owner,
		&svc.TemplateID,
		&svc.WorkspacePath,
		&svc.SourceType,
		&svc.RepoURL,
		&svc.Branch,
		&svc.DockerfilePath,
		&svc.WebhookSecret,
		&svc.AutoDeployEnabled,
		&svc.AutoDeployEnvironment,
		&svc.CreatedAt,
		&svc.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return model.Service{}, ErrNotFound
	}
	if err != nil {
		return model.Service{}, fmt.Errorf("Get service: %w", err)
	}
	return svc, nil
}

func (s *ServicesStore) Create(ctx context.Context, req model.CreateServiceRequest) (model.Service, error) {
	const q = `
		INSERT INTO services (
			name, description, owner, template_id, workspace_path,
			source_type, repo_url, branch, dockerfile_path,
			webhook_secret, auto_deploy_enabled, auto_deploy_environment
		) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id, name, description, owner, template_id, workspace_path,
		source_type, repo_url, branch, dockerfile_path,
		webhook_secret, auto_deploy_enabled, auto_deploy_environment,
		created_at, updated_at`

	var svc model.Service
	err := s.db.QueryRowContext(ctx, q,
		req.Name,
		req.Description,
		req.Owner,
		req.TemplateID,
		req.WorkspacePath,
		req.SourceType,
		req.RepoURL,
		req.Branch,
		req.DockerfilePath,
		req.WebhookSecret,
		req.AutoDeployEnabled,
		req.AutoDeployEnvironment,
	).Scan(
		&svc.ID,
		&svc.Name,
		&svc.Description,
		&svc.Owner,
		&svc.TemplateID,
		&svc.WorkspacePath,
		&svc.SourceType,
		&svc.RepoURL,
		&svc.Branch,
		&svc.DockerfilePath,
		&svc.WebhookSecret,
		&svc.AutoDeployEnabled,
		&svc.AutoDeployEnvironment,
		&svc.CreatedAt,
		&svc.UpdatedAt,
	)
	if err != nil {
		if isUniqueViolation(err) {
			return model.Service{}, ErrConflict
		}
		return model.Service{}, fmt.Errorf("Create service: %w", err)
	}
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
			updated_at = NOW()
		WHERE id = $11
		RETURNING id, name, description, owner, template_id, workspace_path,
		source_type, repo_url, branch, dockerfile_path,
		webhook_secret, auto_deploy_enabled, auto_deploy_environment,
		created_at, updated_at`

	var svc model.Service
	err := s.db.QueryRowContext(ctx, q,
		req.Name,
		req.Description,
		req.Owner,
		req.SourceType,
		req.RepoURL,
		req.Branch,
		req.DockerfilePath,
		req.WebhookSecret,
		req.AutoDeployEnabled,
		req.AutoDeployEnvironment,
		id,
	).Scan(
		&svc.ID,
		&svc.Name,
		&svc.Description,
		&svc.Owner,
		&svc.TemplateID,
		&svc.WorkspacePath,
		&svc.SourceType,
		&svc.RepoURL,
		&svc.Branch,
		&svc.DockerfilePath,
		&svc.WebhookSecret,
		&svc.AutoDeployEnabled,
		&svc.AutoDeployEnvironment,
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
		return model.Service{}, fmt.Errorf("Update service: %w", err)
	}
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
