package repository

import "github.com/jmoiron/sqlx"

type repository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *repository {
	return &repository{db: db}
}

func (r *repository) AddTask(userID int64, task string) error {
	_, err := r.db.Exec("INSERT INTO tasks (user_id, task) VALUES ($1, $2)", userID, task)
	return err
}
