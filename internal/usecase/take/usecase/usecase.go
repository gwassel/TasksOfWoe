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
		Name:      "take",
		Aliases:   nil,
		DescShort: "start working on task(s)",
		DescFull:  "start working on task(s)",
		Format:    "take <ids>",
		Args:      []string{"ids: IDs of tasks to take in work"},
	}
	return &Usecase{logger: logger, taskRepo: taskRepo, Desc: desc}
}

func (u *Usecase) Handle(ctx context.Context, userID int64, userTaskIDs []int64) error {
	err := u.taskRepo.TakeTask(ctx, userID, userTaskIDs)
	if err != nil {
		return err
	}

	return nil
}
