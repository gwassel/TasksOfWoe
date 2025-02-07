package repository

import (
	"github.com/gwassel/TasksOfWoe/internal/domain"
	"github.com/jmoiron/sqlx"
)

type repository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *repository {
	return &repository{db: db}
}

func (r *repository) ListAllTasks(userID int64) ([]domain.Task, error) {
	var tasks []domain.Task
	err := r.db.Select(&tasks, "SELECT * FROM tasks WHERE user_id = $1 ORDER BY completed_at NULLS FIRST, id ASC", userID)
	return tasks, err
}
