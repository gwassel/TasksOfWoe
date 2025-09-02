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
		Name:      "complete",
		Aliases:   []string{"com"},
		DescShort: "complete task(s)",
		DescFull:  "",
		Format:    "desc <ids>",
		Args:      []string{"ids: list of IDs of tasks to mark completed"},
	}
	return &Usecase{logger: logger, taskRepo: taskRepo, Desc: desc}
}

func (u *Usecase) Handle(userID int64, userTaskIDs []int64) error {
	err := u.taskRepo.CompleteTask(userID, userTaskIDs)
	if err != nil {
		return err
	}

	return nil
}
