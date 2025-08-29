package usecase

import domain "github.com/gwassel/TasksOfWoe/internal/domain/task"

type TaskRepo interface {
	GetUnencryptedTasks() ([]domain.Task, error)
	EncryptTask(taskIDs int64, task []byte) error
}
