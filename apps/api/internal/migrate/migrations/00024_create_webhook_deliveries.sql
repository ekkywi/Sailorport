-- +goose Up
CREATE TABLE IF NOT EXISTS webhook_deliveries (
    delivery_id     TEXT PRIMARY KEY,
    service_id      UUID REFERENCES services(id) ON DELETE SET NULL,
    deployment_id   UUID REFERENCES deployments(id) ON DELETE SET NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_webhook_deliveries_created_at
    ON webhook_deliveries(created_at DESC);

-- +goose Down
DROP TABLE IF EXISTS webhook_deliveries;