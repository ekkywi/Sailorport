-- +goose Up
UPDATE runtime_jobs r
SET
    status = 'failed',
    error_message = 'deployment no longer exists',
    updated_at = NOW()
WHERE r.status IN ('pending', 'claimed')
  AND NOT EXISTS (SELECT 1 FROM deployments d WHERE d.id = r.deployment_id);

DELETE FROM runtime_jobs r
WHERE NOT EXISTS (SELECT 1 FROM deployments d WHERE d.id = r.deployment_id);

-- +goose StatementBegin
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'runtime_jobs_deployment_id_fkey'
    ) THEN
        ALTER TABLE runtime_jobs
            ADD CONSTRAINT runtime_jobs_deployment_id_fkey
            FOREIGN KEY (deployment_id) REFERENCES deployments(id) ON DELETE CASCADE;
    END IF;
END $$;
-- +goose StatementEnd

-- +goose Down
ALTER TABLE runtime_jobs
    DROP CONSTRAINT IF EXISTS runtime_jobs_deployment_id_fkey;
