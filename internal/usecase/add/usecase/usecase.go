package usecase

import "github.com/pkg/errors"

type Usecase struct {
	taskRepo TaskRepo
}

func New(taskRepo TaskRepo) *Usecase {
	return &Usecase{taskRepo: taskRepo}
}

func (u *Usecase) Handle(userID int64, task string) (int64, error) {
	userTaskID, err := u.taskRepo.AddTask(userID, task)
	if err != nil {
		return 0, errors.Wrap(err, "failed to add task")
	}
	return userTaskID, nil
}
