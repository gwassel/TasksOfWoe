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
		Name:      "untake",
		Aliases:   nil,
		DescShort: "stop working on task(s)",
		DescFull:  "",
		Format:    "take <ids>",
		Args:      []string{"ids: ids of tasks to stop working on"},
	}
	return &Usecase{logger: logger, taskRepo: taskRepo, Desc: desc}
}

func (u *Usecase) Handle(userID int64, userTaskIDs []int64) error {
	err := u.taskRepo.UntakeTask(userID, userTaskIDs)
	if err != nil {
		return err
	}

	return nil
}
