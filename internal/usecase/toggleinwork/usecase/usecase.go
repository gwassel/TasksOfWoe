package usecase

import (
	"github.com/gwassel/TasksOfWoe/internal/infra"
)

type Usecase struct {
	logger   infra.Logger
	taskRepo TaskRepo
	value    bool
}

func New(logger infra.Logger, taskRepo TaskRepo, value bool) *Usecase {
	return &Usecase{logger: logger, taskRepo: taskRepo, value: value}
}

func (u *Usecase) Handle(userID int64, userTaskIDs []int64) error {
	err := u.taskRepo.ToggleInWork(userID, userTaskIDs, u.value)
	if err != nil {
		return err
	}

	return nil
}
