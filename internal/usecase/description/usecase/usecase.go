package usecase

import (
	"github.com/gwassel/TasksOfWoe/internal/domain"
	"github.com/gwassel/TasksOfWoe/internal/infra"
)

type Usecase struct {
	logger   infra.Logger
	taskRepo TaskRepo
	Desc     domain.Description
}

func New(logger infra.Logger, taskRepo TaskRepo) *Usecase {
	desc := domain.Description{
		Name:      "description",
		Aliases:   []string{"desc"},
		DescShort: "print task description(s)",
		DescFull:  "print task description(s)",
		Format:    "decription <ids>",
		Args:      []string{"ids: IDs of tasks to get description for"},
	}
	return &Usecase{logger: logger, taskRepo: taskRepo, Desc: desc}
}

func (u *Usecase) Handle(userID int64, userTaskIDs []int64) ([]domain.Task, error) {
	tasks, err := u.taskRepo.TaskDescription(userID, userTaskIDs)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}
