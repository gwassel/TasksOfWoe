BEGIN;

ALTER TABLE tasks
    ADD COLUMN user_task_id BIGINT NULL;

UPDATE tasks SET user_task_id=id;

ALTER TABLE tasks
    ALTER COLUMN user_task_id BIGINT NOT NULL;

COMMIT;
