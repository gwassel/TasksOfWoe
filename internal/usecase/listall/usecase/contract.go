package usecase

import domain "github.com/gwassel/TasksOfWoe/internal/domain/task"

type TaskRepo interface {
	ListAllTasks(userID int64) ([]domain.Task, error)
}
