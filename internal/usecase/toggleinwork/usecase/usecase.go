package usecase

import (
	"github.com/gwassel/TasksOfWoe/internal/infra"
	"github.com/pkg/errors"
)

type Usecase struct {
	logger   infra.Logger
	taskRepo TaskRepo
}

func New(logger infra.Logger, taskRepo TaskRepo) *Usecase {
	return &Usecase{logger: logger, taskRepo: taskRepo}
}

func (u *Usecase) Handle(userID int64, userTaskIDs []int64) string {
	for _, taskID := range userTaskIDs {
		// TODO add transaction
		err := u.handleOne(userID, taskID)
		if err != nil {
			u.logger.Error(errors.Wrap(err, "unable to complete task"))
			return "Failed to complete task."
		}
	}

	return "Task status updated"
}

func (u *Usecase) handleOne(userID, userTaskID int64) error {
	err := u.taskRepo.ToggleInWork(userID, userTaskID)
	if err != nil {
		return err
	}
	return nil
}
