-- +goose Up
ALTER TABLE services
    ADD COLUMN IF NOT EXISTS owner_user_id UUID REFERENCES users(id) ON DELETE SET NULL;

CREATE INDEX IF NOT EXISTS idx_services_owner_user_id
    ON services(owner_user_id);

-- Backfill
UPDATE services s
SET owner_user_id = u.id
FROM users u
WHERE s.owner_user_id IS NULL
    AND u.deleted_at IS NULL
    AND lower(trim(s.owner)) = lower(u.email);

-- +goose Down
DROP INDEX IF EXISTS idx_services_owner_user_id;

ALTER TABLE services
    DROP COLUMN IF EXISTS owner_user_id;