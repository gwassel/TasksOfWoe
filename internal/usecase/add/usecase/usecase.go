package usecase

import (
	"os"

	"github.com/gwassel/TasksOfWoe/internal/domain/encoder"
	"github.com/pkg/errors"
)

type Usecase struct {
	taskRepo TaskRepo
}

func New(taskRepo TaskRepo) *Usecase {
	return &Usecase{taskRepo: taskRepo}
}

func (u *Usecase) Handle(userID int64, task string) (int64, error) {
	e, err := encoder.New(os.Getenv("ENCRYPTION_KEY"))
	if err != nil {
		return 0, errors.Wrap(err, "failed to add task")
	}
	ct, err := e.Encode(task)
	if err != nil {
		return 0, errors.Wrap(err, "failed to encode task")
	}
	userTaskID, err := u.taskRepo.AddTask(userID, ct)
	if err != nil {
		return 0, errors.Wrap(err, "failed to add task")
	}
	return userTaskID, nil
}
