package usecase

import "github.com/gwassel/TasksOfWoe/internal/domain"

type TaskRepo interface {
	ListAllTasks(userID int64) ([]domain.Task, error)
}
