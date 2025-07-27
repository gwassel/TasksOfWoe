package repository

import "github.com/jmoiron/sqlx"

type repository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *repository {
	return &repository{db: db}
}

func (r *repository) CompleteTask(userID int64, userTaskID int64) error {
	_, err := r.db.Exec("UPDATE tasks SET completed = TRUE, completed_at=NOW(), is_in_work = FALSE WHERE user_task_id = $1 AND user_id = $2 AND completed = FALSE", userTaskID, userID)
	return err
}
