package usecase

import (
	"github.com/gwassel/TasksOfWoe/internal/domain"
	"github.com/pkg/errors"
)

type Usecase struct {
	taskRepo TaskRepo
	Desc     domain.Description
}

func New(taskRepo TaskRepo) *Usecase {
	desc := domain.Description{
		Name: "listall",
	}
	return &Usecase{taskRepo: taskRepo, Desc: desc}
}

func (u *Usecase) Handle(userID int64) ([]domain.Task, error) {
	tasks, err := u.taskRepo.ListAllTasks(userID)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to get all tasks")
	}

	return tasks, nil
}
