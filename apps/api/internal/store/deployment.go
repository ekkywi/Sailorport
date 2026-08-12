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

func (s *DeploymentsStore) Create(ctx context.Context, serviceID string) (model.Deployment, error) {
	const q = `
		INSERT INTO deployments (service_id, status)
		VALUES ($1, 'pending')
		RETURNING id, service_id, worker_id, status, image_tag, container_id, port, error_message, created_at, updated_at`

	return scanDeployment(s.db.QueryRowContext(ctx, q, serviceID))
}

func (s *DeploymentsStore) Get(ctx context.Context, id string) (model.Deployment, error) {
	const q = `
		SELECT id, service_id, worker_id, status, image_tag, container_id, port, error_message, created_at, updated_at
		FROM deployments
		WHERE id = $1`

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
		SELECT id, service_id, worker_id, status, image_tag, container_id, port, error_message, created_at, updated_at
		FROM deployments
		ORDER BY created_at DESC`

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
		SELECT id, service_id, worker_id, status, image_tag, container_id, port, error_message, created_at, updated_at
		FROM deployments
		WHERE service_id = $1
		ORDER BY created_at DESC`

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
		FROM next_job, services s
		WHERE d.id = next_job.id AND s.id = d.service_id
		RETURNING
			d.id, d.service_id, d.worker_id, d.status, d.image_tag, d.container_id,
			d.port, d.error_message, d.created_at, d.updated_at,
			s.name, s.workspace_path`
	return scanDeploymentJob(s.db.QueryRowContext(ctx, claimQ, workerID))
}

func (s *DeploymentsStore) Update(ctx context.Context, id string, req model.UpdateDeploymentRequest) (model.Deployment, error) {
	const q = `
		UPDATE deployments
		SET
			status = COALESCE(NULLIF($2, ''), status),
			image_tag = COALESCE(NULLIF($3, ''), image_tag),
			container_id = COALESCE(NULLIF($4, ''), container_id),
			port = COALESCE($5, port),
			error_message = COALESCE(NULLIF($6, ''), error_message),
			updated_at = NOW()
		WHERE id = $1
		RETURNING id, service_id, worker_id, status, image_tag, container_id, port, error_message, created_at, updated_at`

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
		SELECT DISTINCT ON (service_id)
			id, service_id, worker_id, status, image_tag, container_id, port, error_message, created_at, updated_at
		FROM deployments
		ORDER BY service_id, created_at DESC
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
