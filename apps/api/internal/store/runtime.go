package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ekkywi/sailorport/apps/api/internal/model"
)

type RuntimeStore struct {
	db *sql.DB
}

func NewRuntimeStore(db *sql.DB) *RuntimeStore {
	return &RuntimeStore{db: db}
}

func (s *RuntimeStore) HasActiveJob(ctx context.Context, deploymentID string) (bool, error) {
	const q = `
		SELECT EXISTS (
			SELECT 1 FROM runtime_jobs
			WHERE deployment_id = $1 
				AND status IN ('pending', 'claimed')
				AND action IN ('stop', 'start', 'remove')
		)`
	var exists bool
	if err := s.db.QueryRowContext(ctx, q, deploymentID).Scan(&exists); err != nil {
		return false, fmt.Errorf("HasActiveJob: %w", err)
	}
	return exists, nil
}

func (s *RuntimeStore) Create(ctx context.Context, serviceID, deploymentID, serviceName, action string) (model.RuntimeJob, error) {
	const q = `
		WITH inserted AS (
			INSERT INTO runtime_jobs (service_id, deployment_id, service_name, action, status)
			VALUES ($1, $2, $3, $4, 'pending')
			RETURNING *
		)
		SELECT
			i.id, i.service_id, i.deployment_id, i.service_name, e.slug, i.action, i.status, i.worker_id, i.error_message, i.output, i.created_at, i.updated_at
		FROM inserted i
		JOIN deployments d ON d.id = i.deployment_id
		JOIN environments e ON e.id = d.environment_id`

	return scanRuntimeJob(s.db.QueryRowContext(ctx, q, serviceID, deploymentID, serviceName, action))
}

func (s *RuntimeStore) ClaimNext(ctx context.Context, workerID string) (model.RuntimeJob, error) {
	const q = `
		WITH next_job AS (
			SELECT r.id, e.slug AS environment_slug
			FROM runtime_jobs r
			INNER JOIN deployments d ON d.id = r.deployment_id
			INNER JOIN environments e ON e.id = d.environment_id
			WHERE r.status = 'pending'
			ORDER BY r.created_at ASC
			LIMIT 1
			FOR UPDATE OF r SKIP LOCKED
		)
		UPDATE runtime_jobs r
		SET status = 'claimed', worker_id = $1::uuid, updated_at = NOW()
		FROM next_job nj
		WHERE r.id = nj.id
		RETURNING
			r.id, r.service_id, r.deployment_id, r.service_name, nj.environment_slug, r.action, r.status,
			r.worker_id, r.error_message, r.output, r.created_at, r.updated_at`

	job, err := scanRuntimeJob(s.db.QueryRowContext(ctx, q, workerID))
	if errors.Is(err, sql.ErrNoRows) {
		return model.RuntimeJob{}, ErrNotFound
	}
	if err != nil {
		return model.RuntimeJob{}, fmt.Errorf("ClaimNext runtime: %w", err)
	}
	return job, nil
}

func (s *RuntimeStore) Update(ctx context.Context, id string, req model.UpdateRuntimeJobRequest) (model.RuntimeJob, error) {
	const q = `
		WITH updated AS (
			UPDATE runtime_jobs
			SET
				status = COALESCE(NULLIF($2, ''), status),
				error_message = COALESCE(NULLIF($3, ''), error_message),
				output = COALESCE(NULLIF($4, ''), output),
				updated_at = NOW()
			WHERE id = $1
			RETURNING *
		)
		SELECT
			u.id, u.service_id, u.deployment_id, u.service_name, e.slug,
			u.action, u.status, u.worker_id, u.error_message, u.output, u.created_at, u.updated_at
		FROM updated u
		JOIN deployments d ON d.id = u.deployment_id
		JOIN environments e ON e.id = d.environment_id
	`

	job, err := scanRuntimeJob(s.db.QueryRowContext(ctx, q, id, req.Status, req.ErrorMessage, req.Output))
	if errors.Is(err, sql.ErrNoRows) {
		return model.RuntimeJob{}, ErrNotFound
	}
	if err != nil {
		return model.RuntimeJob{}, fmt.Errorf("Update runtime job: %w", err)
	}
	return job, nil
}

func (s *RuntimeStore) Get(ctx context.Context, id string) (model.RuntimeJob, error) {
	const q = `
		SELECT
			r.id, r.service_id, r.deployment_id, r.service_name, e.slug,
			r.action, r.status, r.worker_id, r.error_message, r.output,
			r.created_at, r.updated_at
		FROM runtime_jobs r
		JOIN deployments d ON d.id = r.deployment_id
		JOIN environments e ON e.id = d.environment_id
		WHERE r.id = $1`
	job, err := scanRuntimeJob(s.db.QueryRowContext(ctx, q, id))
	if errors.Is(err, sql.ErrNoRows) {
		return model.RuntimeJob{}, ErrNotFound
	}
	if err != nil {
		return model.RuntimeJob{}, fmt.Errorf("Get runtime job: %w", err)
	}
	return job, nil
}

func scanRuntimeJob(row rowScanner) (model.RuntimeJob, error) {
	var (
		job      model.RuntimeJob
		workerID sql.NullString
	)
	err := row.Scan(
		&job.ID,
		&job.ServiceID,
		&job.DeploymentID,
		&job.ServiceName,
		&job.EnvironmentSlug,
		&job.Action,
		&job.Status,
		&workerID,
		&job.ErrorMessage,
		&job.Output,
		&job.CreatedAt,
		&job.UpdatedAt,
	)
	if err != nil {
		return model.RuntimeJob{}, err
	}
	if workerID.Valid {
		v := workerID.String
		job.WorkerID = &v
	}
	return job, nil
}
