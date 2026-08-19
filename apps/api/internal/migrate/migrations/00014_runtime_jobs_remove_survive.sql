-- +goose Up
ALTER TABLE runtime_jobs
    ADD COLUMN IF NOT EXISTS environment_slug TEXT NOT NULL DEFAULT '';

UPDATE runtime_jobs r
SET environment_slug = e.slug
FROM deployments d
JOIN environments e ON e.id = d.environment_id
WHERE r.deployment_id = d.id
  AND r.environment_slug = '';

UPDATE runtime_jobs
SET environment_slug = 'dev'
WHERE environment_slug = '';

ALTER TABLE runtime_jobs
    ALTER COLUMN deployment_id DROP NOT NULL;

ALTER TABLE runtime_jobs
    DROP CONSTRAINT IF EXISTS runtime_jobs_deployment_id_fkey;

ALTER TABLE runtime_jobs
    ADD CONSTRAINT runtime_jobs_deployment_id_fkey
    FOREIGN KEY (deployment_id) REFERENCES deployments(id) ON DELETE SET NULL;

-- +goose Down
DELETE FROM runtime_jobs WHERE deployment_id IS NULL;

ALTER TABLE runtime_jobs
    DROP CONSTRAINT IF EXISTS runtime_jobs_deployment_id_fkey;

ALTER TABLE runtime_jobs
    ALTER COLUMN deployment_id SET NOT NULL;

ALTER TABLE runtime_jobs
    ADD CONSTRAINT runtime_jobs_deployment_id_fkey
    FOREIGN KEY (deployment_id) REFERENCES deployments(id) ON DELETE CASCADE;

ALTER TABLE runtime_jobs
    DROP COLUMN IF EXISTS environment_slug;
