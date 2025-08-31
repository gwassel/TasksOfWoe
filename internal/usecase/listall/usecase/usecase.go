package usecase

import (
	"os"

	"github.com/gwassel/TasksOfWoe/internal/domain/encoder"
	domain "github.com/gwassel/TasksOfWoe/internal/domain/task"
	"github.com/pkg/errors"
)

type Usecase struct {
	taskRepo TaskRepo
}

func New(taskRepo TaskRepo) *Usecase {
	return &Usecase{taskRepo: taskRepo}
}

func (u *Usecase) Handle(userID int64) ([]domain.Task, error) {
	tasks, err := u.taskRepo.ListAllTasks(userID)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to get all tasks")
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
