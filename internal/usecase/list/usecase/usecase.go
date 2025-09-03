package usecase

import (
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
		Name:      "list",
		Aliases:   []string{"ls"},
		DescShort: "list tasks",
		DescFull:  "list uncompleted tasks",
		Format:    "list",
		Args:      nil,
	}
	return &Usecase{logger: logger, taskRepo: taskRepo, Desc: desc}
}

func (u *Usecase) Handle(userID int64) ([]domain.Task, error) {
	tasks, err := u.taskRepo.ListTasks(userID)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}
