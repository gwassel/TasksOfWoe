BEGIN;

ALTER TABLE tasks
    ADD CONSTRAINT tasks_unique_user_task_id_user__ctr UNIQUE (user_id, user_task_id);

COMMIT;
