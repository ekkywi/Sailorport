-- +goose Up
ALTER TABLE deployments
    ADD COLUMN IF NOT EXISTS target_worker_id UUID REFERENCES workers(id) ON DELETE SET NULL;

ALTER TABLE runtime_jobs
    ADD COLUMN IF NOT EXISTS target_worker_id UUID REFERENCES workers(id) ON DELETE SET NULL;

CREATE INDEX IF NOT EXISTS idx_deployments_target_worker_pending
    ON deployments (target_worker_id)
    WHERE status = 'pending';

CREATE INDEX IF NOT EXISTS idx_runtime_jobs_target_worker_pending
    ON runtime_jobs (target_worker_id)
    WHERE status = 'pending';

-- +goose Down
DROP INDEX IF EXISTS idx_runtime_jobs_target_worker_pending;
DROP INDEX IF EXISTS idx_deployments_target_worker_pending;

ALTER TABLE runtime_jobs
    DROP COLUMN IF EXISTS target_worker_id;

ALTER TABLE deployments
    DROP COLUMN IF EXISTS target_worker_id;