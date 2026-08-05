-- +goose Up
CREATE TABLE IF NOT EXISTS services (
    id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    name        TEXT        NOT NULL    UNIQUE,
    description TEXT        NOT NULL    DEFAULT '',
    owner       TEXT        NOT NULL    DEFAULT '',
    created_at  TIMESTAMPTZ NOT NULL    DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL    DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS services;