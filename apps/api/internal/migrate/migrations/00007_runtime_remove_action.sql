-- +goose Up
ALTER TABLE runtime_jobs
    DROP CONSTRAINT IF EXISTS runtime_jobs_action_check;

ALTER TABLE runtime_jobs
    ADD CONSTRAINT runtime_jobs_action_check
    CHECK (action IN ('stop', 'start', 'remove'));

-- +goose Down
ALTER TABLE runtime_jobs
    DROP CONSTRAINT IF EXISTS runtime_jobs_action_check;
    
ALTER TABLE runtime_jobs
    ADD CONSTRAINT runtime_jobs_action_check
    CHECK (action IN ('stop', 'start'));