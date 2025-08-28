package usecase

import domain "github.com/gwassel/TasksOfWoe/internal/domain/task"

type TaskRepo interface {
	TaskDescription(userID int64, userTaskIDs []int64) ([]domain.Task, error)
}
