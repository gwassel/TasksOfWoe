package repository

import "github.com/jmoiron/sqlx"

type repository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *repository {
	return &repository{db: db}
}

func (r *repository) ToggleInWork(userID int64, userTaskID int64) error {
	_, err := r.db.Exec("UPDATE tasks SET is_in_work = NOT is_in_work WHERE user_task_id = $1 AND user_id = $2", userTaskID, userID)
	return err
}
