-- +goose Up
ALTER TABLE services
    ADD COLUMN IF NOT EXISTS source_type TEXT NOT NULL DEFAULT 'scaffold',
    ADD COLUMN IF NOT EXISTS repo_url TEXT NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS branch TEXT NOT NULL DEFAULT 'main',
    ADD COLUMN IF NOT EXISTS dockerfile_path TEXT NOT NULL DEFAULT 'Dockerfile';

-- +goose Down
ALTER TABLE services
    DROP COLUMN IF EXISTS dockerfile_path,
    DROP COLUMN IF EXISTS branch,
    DROP COLUMN IF EXISTS repo_url,
    DROP COLUMN IF EXISTS source_type;