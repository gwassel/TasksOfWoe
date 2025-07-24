package repository

import "github.com/jmoiron/sqlx"

type repository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *repository {
	return &repository{db: db}
}

func (r *repository) AddTask(userID int64, task string) error {
	_, err := r.db.Exec("INSERT INTO tasks (user_id, task, user_task_id) VALUES ($1, $2, (SELECT MAX(user_task_id)+1 FROM tasks where user_id=$3))", userID, task, userID)
	return err
}
