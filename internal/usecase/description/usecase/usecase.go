package usecase

import (
	"os"

	"github.com/gwassel/TasksOfWoe/internal/domain/encoder"
	domain "github.com/gwassel/TasksOfWoe/internal/domain/task"
	"github.com/gwassel/TasksOfWoe/internal/infra"
)

type Usecase struct {
	logger   infra.Logger
	taskRepo TaskRepo
}

func New(logger infra.Logger, taskRepo TaskRepo) *Usecase {
	return &Usecase{logger: logger, taskRepo: taskRepo}
}

func (u *Usecase) Handle(userID int64, userTaskIDs []int64) ([]domain.Task, error) {
	tasks, err := u.taskRepo.TaskDescription(userID, userTaskIDs)
	if err != nil {
		return nil, err
	}

	e, err := encoder.New(os.Getenv("ENCRYPTION_KEY"))
	if err != nil {
		return nil, err
	}
	for i, task := range tasks {
		tasks[i].Task, err = e.Decode(task.EncryptedTask)
		if err != nil {
			return nil, err
		}
	}

	return tasks, nil
}
