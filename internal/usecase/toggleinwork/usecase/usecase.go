package usecase

import (
	"github.com/gwassel/TasksOfWoe/internal/infra"
)

type Usecase struct {
	logger   infra.Logger
	taskRepo TaskRepo
}

func New(logger infra.Logger, taskRepo TaskRepo) *Usecase {
	return &Usecase{logger: logger, taskRepo: taskRepo}
}

func (u *Usecase) Handle(userID int64, userTaskIDs []int64, value bool) error {
	for _, taskID := range userTaskIDs {
		// TODO add transaction
		err := u.handleOne(userID, taskID, value)
		if err != nil {
			return err
		}
	}

	return nil
}

func (u *Usecase) handleOne(userID, userTaskID int64, value bool) error {
	err := u.taskRepo.ToggleInWork(userID, userTaskID, value)
	if err != nil {
		return err
	}
	return nil
}
