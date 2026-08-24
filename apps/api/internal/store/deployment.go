package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ekkywi/sailorport/apps/api/internal/model"
)

type DeploymentsStore struct {
	db *sql.DB
}

func NewDeploymentsStore(db *sql.DB) *DeploymentsStore {
	return &DeploymentsStore{db: db}
}

func (s *DeploymentsStore) Create(ctx context.Context, serviceID, environmentID string, targetWorkerID *string, gitSHA string) (model.Deployment, error) {
	const q = `
		WITH inserted AS (
			INSERT INTO deployments (service_id, environment_id, status, target_worker_id, git_sha)
			VALUES ($1, $2, 'pending', $3, $4)
			RETURNING *
		)
		SELECT
			i.id, i.service_id, i.environment_id, e.slug,
			i.target_worker_id, i.worker_id, i.status, i.image_tag, i.git_sha, i.container_id, i.port,
			i.error_message, i.created_at, i.updated_at
		FROM inserted i
		JOIN environments e ON e.id = i.environment_id`

	return scanDeployment(s.db.QueryRowContext(ctx, q, serviceID, environmentID, targetWorkerID, gitSHA))
}

func (s *DeploymentsStore) Get(ctx context.Context, id string) (model.Deployment, error) {
	const q = `
		SELECT
			d.id, d.service_id, d.environment_id, e.slug,
			d.target_worker_id, d.worker_id, d.status, d.image_tag, d.git_sha, d.container_id, d.port,
			d.error_message, d.created_at, d.updated_at
		FROM deployments d
		JOIN environments e ON e.id = d.environment_id
		WHERE d.id = $1`

	d, err := scanDeployment(s.db.QueryRowContext(ctx, q, id))
	if errors.Is(err, sql.ErrNoRows) {
		return model.Deployment{}, ErrNotFound
	}
	if err != nil {
		return model.Deployment{}, fmt.Errorf("Get deployment: %w", err)
	}
	return d, nil
}

func (s *DeploymentsStore) List(ctx context.Context) ([]model.Deployment, error) {
	const q = `
		SELECT
			d.id, d.service_id, d.environment_id, e.slug,
			d.target_worker_id, d.worker_id, d.status, d.image_tag, d.git_sha, d.container_id, d.port,
			d.error_message, d.created_at, d.updated_at
		FROM deployments d
		JOIN environments e ON e.id = d.environment_id
		ORDER BY d.created_at DESC`

	rows, err := s.db.QueryContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("List deployment: %w", err)
	}
	defer rows.Close()

	out := make([]model.Deployment, 0)
	for rows.Next() {
		d, err := scanDeployment(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, d)
	}
	return out, rows.Err()
}

func (s *DeploymentsStore) ListByService(ctx context.Context, serviceID string) ([]model.Deployment, error) {
	const q = `
		SELECT
			d.id, d.service_id, d.environment_id, e.slug,
			d.target_worker_id, d.worker_id, d.status, d.image_tag, d.git_sha, d.container_id, d.port,
			d.error_message, d.created_at, d.updated_at
		FROM deployments d
		JOIN environments e ON e.id = d.environment_id
		WHERE d.service_id = $1
		ORDER BY d.created_at DESC`

	rows, err := s.db.QueryContext(ctx, q, serviceID)
	if err != nil {
		return nil, fmt.Errorf("List deployments by service: %w", err)
	}
	defer rows.Close()

	out := make([]model.Deployment, 0)
	for rows.Next() {
		d, err := scanDeployment(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, d)
	}
	return out, rows.Err()
}

func (s *DeploymentsStore) ClaimNext(ctx context.Context, workerID string) (model.DeploymentJob, error) {
	const claimQ = `
		WITH next_job AS (
			SELECT id FROM deployments
			WHERE status = 'pending'
			  AND (target_worker_id IS NULL OR target_worker_id = $1::uuid)
			ORDER BY created_at ASC
			LIMIT 1
			FOR UPDATE SKIP LOCKED
		)
		UPDATE deployments d
		SET status = 'claimed', worker_id = $1::uuid, updated_at = NOW()
		FROM next_job, services s, environments e
		WHERE d.id = next_job.id
			AND s.id = d.service_id
			AND e.id = d.environment_id
		RETURNING
			d.id, d.service_id, d.environment_id, e.slug,
			d.target_worker_id, d.worker_id, d.status, d.image_tag, d.git_sha, d.container_id,
			d.port, d.error_message, d.created_at, d.updated_at,
			s.name, s.workspace_path,
			s.source_type, s.repo_url, s.branch, s.dockerfile_path`
	return scanDeploymentJob(s.db.QueryRowContext(ctx, claimQ, workerID))
}

func (s *DeploymentsStore) Update(ctx context.Context, id string, req model.UpdateDeploymentRequest) (model.Deployment, error) {
	const q = `
		WITH updated AS (
			UPDATE deployments
			SET
				status = COALESCE(NULLIF($2, ''), status),
				image_tag = COALESCE(NULLIF($3, ''), image_tag),
				git_sha = COALESCE(NULLIF($4, ''), git_sha),
				container_id = COALESCE(NULLIF($5, ''), container_id),
				port = COALESCE($6, port),
				error_message = COALESCE(NULLIF($7, ''), error_message),
				updated_at = NOW()
			WHERE id = $1
			RETURNING *
		)
		SELECT
			u.id, u.service_id, u.environment_id, e.slug,
			u.target_worker_id, u.worker_id, u.status, u.image_tag, u.git_sha, u.container_id, u.port,
			u.error_message, u.created_at, u.updated_at
		FROM updated u
		JOIN environments e ON e.id = u.environment_id`

	d, err := scanDeployment(s.db.QueryRowContext(ctx, q, id, req.Status, req.ImageTag, req.GitSHA, req.ContainerID, req.Port, req.ErrorMessage))
	if errors.Is(err, sql.ErrNoRows) {
		return model.Deployment{}, ErrNotFound
	}
	if err != nil {
		return model.Deployment{}, fmt.Errorf("Update deployment: %w", err)
	}
	return d, nil
}

func scanDeployment(row rowScanner) (model.Deployment, error) {
	var (
		d              model.Deployment
		targetWorkerID sql.NullString
		workerID       sql.NullString
		port           sql.NullInt64
	)
	err := row.Scan(
		&d.ID,
		&d.ServiceID,
		&d.EnvironmentID,
		&d.EnvironmentSlug,
		&targetWorkerID,
		&workerID,
		&d.Status,
		&d.ImageTag,
		&d.GitSHA,
		&d.ContainerID,
		&port,
		&d.ErrorMessage,
		&d.CreatedAt,
		&d.UpdatedAt,
	)
	if err != nil {
		return model.Deployment{}, err
	}
	if targetWorkerID.Valid {
		v := targetWorkerID.String
		d.TargetWorkerID = &v
	}
	if workerID.Valid {
		v := workerID.String
		d.WorkerID = &v
	}
	if port.Valid {
		p := int(port.Int64)
		d.Port = &p
	}
	return d, nil
}

func scanDeploymentJob(row rowScanner) (model.DeploymentJob, error) {
	var (
		job            model.DeploymentJob
		targetWorkerID sql.NullString
		workerID       sql.NullString
		port           sql.NullInt64
	)
	err := row.Scan(
		&job.ID,
		&job.ServiceID,
		&job.EnvironmentID,
		&job.EnvironmentSlug,
		&targetWorkerID,
		&workerID,
		&job.Status,
		&job.ImageTag,
		&job.GitSHA,
		&job.ContainerID,
		&port,
		&job.ErrorMessage,
		&job.CreatedAt,
		&job.UpdatedAt,
		&job.ServiceName,
		&job.WorkspacePath,
		&job.SourceType,
		&job.RepoURL,
		&job.Branch,
		&job.DockerfilePath,
	)
	if err != nil {
		return model.DeploymentJob{}, err
	}
	if targetWorkerID.Valid {
		v := targetWorkerID.String
		job.TargetWorkerID = &v
	}
	if workerID.Valid {
		v := workerID.String
		job.WorkerID = &v
	}
	if port.Valid {
		p := int(port.Int64)
		job.Port = &p
	}
	return job, nil
}

func (s *DeploymentsStore) LatestByServices(ctx context.Context) (map[string]model.Deployment, error) {
	const q = `
		SELECT DISTINCT ON (d.service_id)
			d.id, d.service_id, d.environment_id, e.slug,
			d.target_worker_id, d.worker_id, d.status, d.image_tag, d.git_sha, d.container_id, d.port,
			d.error_message, d.created_at, d.updated_at
		FROM deployments d
		JOIN environments e ON e.id = d.environment_id
		ORDER BY d.service_id, d.created_at DESC
	`

	rows, err := s.db.QueryContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("LatestByServices: %w", err)
	}
	defer rows.Close()

	out := make(map[string]model.Deployment)
	for rows.Next() {
		d, err := scanDeployment(rows)
		if err != nil {
			return nil, err
		}
		out[d.ServiceID] = d
	}
	return out, rows.Err()
}

// LatestPerEnvByServices returns the newest deployment per (service, environment).
// Outer map key = service_id, inner map key = environment slug (dev, staging, prod).
func (s *DeploymentsStore) LatestPerEnvByServices(ctx context.Context) (map[string]map[string]model.Deployment, error) {
	const q = `
		SELECT DISTINCT ON (d.service_id, d.environment_id)
			d.id, d.service_id, d.environment_id, e.slug,
			d.target_worker_id, d.worker_id, d.status, d.image_tag, d.git_sha, d.container_id, d.port,
			d.error_message, d.created_at, d.updated_at
		FROM deployments d
		JOIN environments e ON e.id = d.environment_id
		ORDER BY d.service_id, d.environment_id, d.created_at DESC
	`

	rows, err := s.db.QueryContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("LatestPerEnvByServices: %w", err)
	}
	defer rows.Close()
	out := make(map[string]map[string]model.Deployment)
	for rows.Next() {
		d, err := scanDeployment(rows)
		if err != nil {
			return nil, err
		}
		if out[d.ServiceID] == nil {
			out[d.ServiceID] = make(map[string]model.Deployment)
		}
		out[d.ServiceID][d.EnvironmentSlug] = d
	}
	return out, rows.Err()
}
