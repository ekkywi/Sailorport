-- +goose Up
ALTER TABLE services
    ADD COLUMN IF NOT EXISTS template_id TEXT NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS workspace_path TEXT NOT NULL DEFAULT '';

-- +goose Down
ALTER TABLE services
    DROP COLUMN IF EXISTS template_id,
    DROP COLUMN IF EXISTS workspace_path;