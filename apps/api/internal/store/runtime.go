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

func (s *RuntimeStore) Create(ctx context.Context, serviceID, deploymentID, serviceName, action string) (model.RuntimeJob, error) {
	const q = `
		INSERT INTO runtime_jobs (service_id, deployment_id, service_name, action, status)
		VALUES ($1, $2, $3, $4, 'pending')
		RETURNING id, service_id, deployment_id, service_name, action, status, worker_id, error_message, created_at, updated_at`

	return scanRuntimeJob(s.db.QueryRowContext(ctx, q, serviceID, deploymentID, serviceName, action))
}

func (s *RuntimeStore) ClaimNext(ctx context.Context, workerID string) (model.RuntimeJob, error) {
	const q = `
		WITH next_job AS (
			SELECT id FROM runtime_jobs
			WHERE status = 'pending'
			ORDER BY created_at ASC
			LIMIT 1
			FOR UPDATE SKIP LOCKED
		)
		UPDATE runtime_jobs r
		SET status = 'claimed', worker_id = $1::uuid, updated_at = NOW()
		FROM next_job
		WHERE r.id = next_job.id
		RETURNING r.id, r.service_id, r.deployment_id, r.service_name, r.action, r.status,
		          r.worker_id, r.error_message, r.created_at, r.updated_at`

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
		UPDATE runtime_jobs
		SET
			status = COALESCE(NULLIF($2, ''), status),
			error_message = COALESCE(NULLIF($3, ''), error_message),
			updated_at = NOW()
		WHERE id = $1
		RETURNING id, service_id, deployment_id, service_name, action, status, worker_id, error_message, created_at, updated_at`

	job, err := scanRuntimeJob(s.db.QueryRowContext(ctx, q, id, req.Status, req.ErrorMessage))
	if errors.Is(err, sql.ErrNoRows) {
		return model.RuntimeJob{}, ErrNotFound
	}
	if err != nil {
		return model.RuntimeJob{}, fmt.Errorf("Update runtime job: %w", err)
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
		&job.Action,
		&job.Status,
		&workerID,
		&job.ErrorMessage,
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
