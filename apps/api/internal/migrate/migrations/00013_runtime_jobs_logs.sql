-- +goose Up
ALTER TABLE runtime_jobs
    DROP CONSTRAINT IF EXISTS runtime_jobs_action_check;

ALTER TABLE runtime_jobs
    ADD CONSTRAINT runtime_jobs_action_check
    CHECK (action IN ('stop', 'start', 'remove', 'logs'));

ALTER TABLE runtime_jobs
    ADD COLUMN IF NOT EXISTS output TEXT NOT NULL DEFAULT '';

-- +goose Down
ALTER TABLE runtime_jobs
    DROP COLUMN IF EXISTS output;

ALTER TABLE runtime_jobs
    DROP CONSTRAINT IF EXISTS runtime_jobs_action_check;

ALTER TABLE runtime_jobs
    ADD CONSTRAINT runtime_jobs_action_check
    CHECK (action IN ('stop', 'start', 'remove'));