package usecase

import "github.com/gwassel/TasksOfWoe/internal/domain/task"

type TaskRepo interface {
	ListTasks(userID int64) ([]domain.Task, error)
}
