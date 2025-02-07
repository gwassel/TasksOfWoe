package usecase

import "github.com/gwassel/TasksOfWoe/internal/domain"

type TaskRepo interface {
	ListTasks(userID int64) ([]domain.Task, error)
}
