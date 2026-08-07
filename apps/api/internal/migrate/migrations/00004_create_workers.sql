-- +goose Up
CREATE TABLE IF NOT EXISTS workers (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name         TEXT NOT NULL,
    hostname     TEXT NOT NULL DEFAULT '',
    status       TEXT NOT NULL DEFAULT 'offline'
                 CHECK (status IN ('online', 'offline', 'draining')),
    labels       JSONB NOT NULL DEFAULT '{}'::jsonb,
    last_seen_at TIMESTAMPTZ,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (name)
);

-- +goose Down
DROP TABLE IF EXISTS workers;