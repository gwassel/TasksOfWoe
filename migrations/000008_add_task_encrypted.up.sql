BEGIN;

ALTER TABLE tasks
    RENAME COLUMN task TO task_unencrypted;

ALTER TABLE tasks
    ADD COLUMN task BYTEA NULL;

COMMIT;
