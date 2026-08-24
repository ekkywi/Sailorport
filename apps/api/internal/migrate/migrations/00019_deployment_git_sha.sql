-- +goose Up
ALTER TABLE deployments
    ADD COLUMN IF NOT EXISTS git_sha TEXT NOT NULL DEFAULT '';

-- +goose Down
ALTER TABLE deployments
    DROP COLUMN IF EXISTS git_sha;