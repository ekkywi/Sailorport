-- +goose Up
CREATE TABLE IF NOT EXISTS app_settings (
    key         TEXT PRIMARY KEY,
    value       TEXT NOT NULL,
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

INSERT INTO app_settings (key, value) VALUES
    ('registration_open', 'false')
ON CONFLICT (key) DO NOTHING;

-- +goose Down
DROP TABLE IF EXISTS app_settings;