package usecase

import (
	"context"

	domain "github.com/gwassel/TasksOfWoe/internal/domain/task"
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
		DescFull:  "mark task(s) as completed",
		Format:    "complete <ids>",
		Args:      []string{"ids: IDs of tasks to complete"},
	}
	return &Usecase{logger: logger, taskRepo: taskRepo, Desc: desc}
}

func (u *Usecase) Handle(ctx context.Context, userID int64, userTaskIDs []int64) error {
	err := u.taskRepo.CompleteTask(ctx, userID, userTaskIDs)
	if err != nil {
		return err
	}

	return nil
}
