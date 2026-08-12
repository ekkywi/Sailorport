-- +goose Up
CREATE TABLE IF NOT EXISTS runtime_jobs (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    service_id      UUID NOT NULL,
    deployment_id   UUID NOT NULL,
    service_name    TEXT NOT NULL,
    action          TEXT NOT NULL CHECK (action IN ('stop', 'start')),
    status          TEXT NOT NULL DEFAULT 'pending'
                    CHECK (status IN ('pending', 'claimed', 'done', 'failed')),
    worker_id       UUID REFERENCES workers(id) ON DELETE SET NULL,
    error_message   TEXT NOT NULL DEFAULT '',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_runtime_jobs_status ON runtime_jobs(status);
CREATE INDEX IF NOT EXISTS idx_runtime_jobs_service ON runtime_jobs(service_id);

-- +goose Down
DROP TABLE IF EXISTS runtime_jobs;
