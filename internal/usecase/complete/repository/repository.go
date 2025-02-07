package repository

import "github.com/jmoiron/sqlx"

type repository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *repository {
	return &repository{db: db}
}

func (r *repository) CompleteTask(userID int64, taskID int64) error {
	_, err := r.db.Exec("UPDATE tasks SET completed = TRUE, completed_at=NOW() WHERE id = $1 AND user_id = $2 AND completed = FALSE", taskID, userID)
	return err
}
