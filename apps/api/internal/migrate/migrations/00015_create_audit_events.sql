-- +goose Up
CREATE TABLE IF NOT EXISTS audit_events (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    at             TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    actor_id       UUID REFERENCES users(id) ON DELETE SET NULL,
    actor_email    TEXT NOT NULL DEFAULT '',
    action         TEXT NOT NULL,
    resource_type  TEXT NOT NULL,
    resource_id    TEXT NOT NULL DEFAULT '',
    resource_name  TEXT NOT NULL DEFAULT '',
    payload        JSONB NOT NULL DEFAULT '{}'::jsonb
);

CREATE INDEX IF NOT EXISTS idx_audit_events_at ON audit_events (at DESC);
CREATE INDEX IF NOT EXISTS idx_audit_events_action ON audit_events (action);
CREATE INDEX IF NOT EXISTS idx_audit_events_resource ON audit_events (resource_type, resource_id);

-- +goose Down
DROP INDEX IF EXISTS idx_audit_events_resource;
DROP INDEX IF EXISTS idx_audit_events_action;
DROP INDEX IF EXISTS idx_audit_events_at;
DROP TABLE IF EXISTS audit_events;