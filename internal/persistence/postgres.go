package persistence

import (
	domain "github.com/gwassel/TasksOfWoe/internal/domain/task"
	"github.com/jmoiron/sqlx"
)

type DB struct {
	*sqlx.DB
}

func NewDB(connStr string) (*DB, error) {
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		return nil, err
	}
	return &DB{db}, nil
}

func (db *DB) AddTask(userID int64, task string) error {
	_, err := db.Exec("INSERT INTO tasks (user_id, task) VALUES ($1, $2)", userID, task)
	return err
}

func (db *DB) ListTasks(userID int64) ([]domain.Task, error) {
	var tasks []domain.Task
	err := db.Select(
		&tasks,
		"SELECT * FROM tasks WHERE user_id = $1 AND completed=FALSE ORDER BY id ASC",
		userID,
	)
	return tasks, err
}

func (db *DB) ListAllTasks(userID int64) ([]domain.Task, error) {
	var tasks []domain.Task
	err := db.Select(
		&tasks,
		"SELECT * FROM tasks WHERE user_id = $1 ORDER BY completed_at NULLS FIRST, id ASC",
		userID,
	)
	return tasks, err
}

func (db *DB) CompleteTask(userID int64, taskID int64) error {
	_, err := db.Exec(
		"UPDATE tasks SET completed = TRUE, completed_at=NOW() WHERE id = $1 AND user_id = $2 AND completed = FALSE",
		taskID,
		userID,
	)
	return err
}
