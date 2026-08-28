-- +goose Up
CREATE TABLE IF NOT EXISTS service_catalog_env (
    service_id UUID NOT NULL REFERENCES services(id) ON DELETE CASCADE,
    key        TEXT NOT NULL,
    value      TEXT NOT NULL,
    secret     BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (service_id, key)
);

CREATE INDEX IF NOT EXISTS idx_service_catalog_env_service_id
    ON service_catalog_env (service_id);

-- +goose Down
DROP TABLE IF EXISTS service_catalog_env;