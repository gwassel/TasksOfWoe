package usecase

import (
	"context"

	domain "github.com/gwassel/TasksOfWoe/internal/domain/task"
	"github.com/pkg/errors"
)

type Usecase struct {
	taskRepo TaskRepo
	Desc     domain.Description
}

func New(taskRepo TaskRepo) *Usecase {
	desc := domain.Description{
		Name:      "add",
		Aliases:   nil,
		DescShort: "add a new task",
		DescFull:  "add a new task with given description",
		Format:    "add <task>",
		Args:      []string{"task: task's description"},
	}
	return &Usecase{taskRepo: taskRepo, Desc: desc}
}

func (u *Usecase) Handle(ctx context.Context, userID int64, task string) (int64, error) {
	userTaskID, err := u.taskRepo.AddTask(ctx, userID, task)
	if err != nil {
		return 0, errors.Wrap(err, "failed to add task")
	}
	return userTaskID, nil
}
