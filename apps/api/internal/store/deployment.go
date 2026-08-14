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

func (s *DeploymentsStore) Create(ctx context.Context, serviceID, environmentID string) (model.Deployment, error) {
	const q = `
		WITH inserted AS (
			INSERT INTO deployments (service_id, environment_id, status)
			VALUES ($1, $2, 'pending')
			RETURNING *
		)
		SELECT
			i.id, i.service_id, i.environment_id, e.slug,
			i.worker_id, i.status, i.image_tag, i.container_id, i.port,
			i.error_message, i.created_at, i.updated_at
		FROM inserted i
		JOIN environments e ON e.id = i.environment_id`

	return scanDeployment(s.db.QueryRowContext(ctx, q, serviceID, environmentID))
}

func (s *DeploymentsStore) Get(ctx context.Context, id string) (model.Deployment, error) {
	const q = `
		SELECT
			d.id, d.service_id, d.environment_id, e.slug,
			d.worker_id, d.status, d.image_tag, d.container_id, d.port,
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
			d.worker_id, d.status, d.image_tag, d.container_id, d.port,
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
			d.worker_id, d.status, d.image_tag, d.container_id, d.port,
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
			d.worker_id, d.status, d.image_tag, d.container_id,
			d.port, d.error_message, d.created_at, d.updated_at,
			s.name, s.workspace_path`
	return scanDeploymentJob(s.db.QueryRowContext(ctx, claimQ, workerID))
}

func (s *DeploymentsStore) Update(ctx context.Context, id string, req model.UpdateDeploymentRequest) (model.Deployment, error) {
	const q = `
		WITH updated AS (
			UPDATE deployments
			SET
				status = COALESCE(NULLIF($2, ''), status),
				image_tag = COALESCE(NULLIF($3, ''), image_tag),
				container_id = COALESCE(NULLIF($4, ''), container_id),
				port = COALESCE($5, port),
				error_message = COALESCE(NULLIF($6, ''), error_message),
				updated_at = NOW()
			WHERE id = $1
			RETURNING *
		)
		SELECT
			u.id, u.service_id, u.environment_id, e.slug,
			u.worker_id, u.status, u.image_tag, u.container_id, u.port,
			u.error_message, u.created_at, u.updated_at
		FROM updated u
		JOIN environments e ON e.id = u.environment_id`

	d, err := scanDeployment(s.db.QueryRowContext(ctx, q, id, req.Status, req.ImageTag, req.ContainerID, req.Port, req.ErrorMessage))
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
		d        model.Deployment
		workerID sql.NullString
		port     sql.NullInt64
	)
	err := row.Scan(
		&d.ID,
		&d.ServiceID,
		&d.EnvironmentID,
		&d.EnvironmentSlug,
		&workerID,
		&d.Status,
		&d.ImageTag,
		&d.ContainerID,
		&port,
		&d.ErrorMessage,
		&d.CreatedAt,
		&d.UpdatedAt,
	)
	if err != nil {
		return model.Deployment{}, err
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
		job      model.DeploymentJob
		workerID sql.NullString
		port     sql.NullInt64
	)
	err := row.Scan(
		&job.ID,
		&job.ServiceID,
		&job.EnvironmentID,
		&job.EnvironmentSlug,
		&workerID,
		&job.Status,
		&job.ImageTag,
		&job.ContainerID,
		&port,
		&job.ErrorMessage,
		&job.CreatedAt,
		&job.UpdatedAt,
		&job.ServiceName,
		&job.WorkspacePath,
	)
	if err != nil {
		return model.DeploymentJob{}, err
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
			d.worker_id, d.status, d.image_tag, d.container_id, d.port,
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
