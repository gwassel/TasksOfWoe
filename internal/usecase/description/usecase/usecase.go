package usecase

import (
	"github.com/gwassel/TasksOfWoe/internal/domain"
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
	var tasks []domain.Task
	for _, taskID := range userTaskIDs {
		// TODO add transaction
		task, err := u.handleOne(userID, taskID)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task[0])
	}

	return tasks, nil
}

func (u *Usecase) handleOne(userID int64, userTaskID int64) ([]domain.Task, error) {
	tasks, err := u.taskRepo.TaskDescription(userID, userTaskID)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}
