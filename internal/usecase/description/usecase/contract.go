package usecase

import "github.com/gwassel/TasksOfWoe/internal/domain"

type TaskRepo interface {
	TaskDescription(userID int64, userTaskID int64) ([]domain.Task, error)
}
