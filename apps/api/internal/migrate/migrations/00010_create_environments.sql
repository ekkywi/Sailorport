-- +goose Up
CREATE TABLE IF NOT EXISTS environments (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    slug       TEXT NOT NULL UNIQUE,
    name       TEXT NOT NULL,
    sort_order INT  NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

INSERT INTO environments (slug, name, sort_order) VALUES
    ('dev',     'Development', 1),
    ('staging', 'Staging',     2),
    ('prod',    'Production',  3)
ON CONFLICT (slug) DO NOTHING;

-- +goose Down
DROP TABLE IF EXISTS environments;