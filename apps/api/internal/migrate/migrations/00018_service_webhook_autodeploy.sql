-- +goose Up
ALTER TABLE services
    ADD COLUMN IF NOT EXISTS webhook_secret TEXT NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS auto_deploy_enabled BOOLEAN NOT NULL DEFAULT false,
    ADD COLUMN IF NOT EXISTS auto_deploy_environment TEXT NOT NULL DEFAULT 'staging';

-- +goose Down
ALTER TABLE services
    DROP COLUMN IF EXISTS auto_deploy_environment,
    DROP COLUMN IF EXISTS auto_deploy_enabled,
    DROP COLUMN IF EXISTS webhook_secret;
