-- +goose Up
CREATE TABLE IF NOT EXISTS deployments (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    service_id      UUID NOT NULL REFERENCES services(id) ON DELETE CASCADE,
    worker_id       UUID REFERENCES workers(id) ON DELETE SET NULL,
    status          TEXT NOT NULL DEFAULT 'pending'
                    CHECK (status IN ('pending', 'claimed', 'building', 'running', 'failed', 'stopped')),
    image_tag       TEXT NOT NULL DEFAULT '',
    container_id    TEXT NOT NULL DEFAULT '',
    port            INT,
    error_message   TEXT NOT NULL DEFAULT '',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_deployments_status ON deployments(status);
CREATE INDEX IF NOT EXISTS idx_deployments_service ON deployments(service_id);

-- +goose Down
DROP TABLE IF EXISTS deployments;