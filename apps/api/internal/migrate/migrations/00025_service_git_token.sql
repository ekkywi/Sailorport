-- +goose Up
ALTER TABLE services
    ADD COLUMN IF NOT EXISTS git_token TEXT NOT NULL DEFAULT '';

-- +goose Down
ALTER TABLE services
    DROP COLUMN IF EXISTS git_token;