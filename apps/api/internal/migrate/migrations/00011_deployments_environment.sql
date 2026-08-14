-- +goose Up
ALTER TABLE deployments
    ADD COLUMN IF NOT EXISTS environment_id UUID REFERENCES environments(id);

UPDATE deployments d
SET environment_id = e.id
FROM environments e
WHERE d.environment_id IS NULL AND e.slug = 'dev';

ALTER TABLE deployments
    ALTER COLUMN environment_id SET NOT NULL;

CREATE INDEX IF NOT EXISTS idx_deployments_environment
    ON deployments(environment_id);

CREATE INDEX IF NOT EXISTS idx_deployments_service_env
    ON deployments(service_id, environment_id);

-- +goose Down
DROP INDEX IF EXISTS idx_deployments_service_env;
DROP INDEX IF EXISTS idx_deployments_environment;
ALTER TABLE deployments DROP COLUMN IF EXISTS environment_id;