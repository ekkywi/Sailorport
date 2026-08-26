-- +goose Up
ALTER TABLE services
    ADD COLUMN IF NOT EXISTS catalog_app_id TEXT NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS image TEXT NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS container_port INT NOT NULL DEFAULT 0;
    
-- +goose Down
ALTER TABLE services
    DROP COLUMN IF EXISTS container_port,
    DROP COLUMN IF EXISTS image,
    DROP COLUMN IF EXISTS catalog_app_id;