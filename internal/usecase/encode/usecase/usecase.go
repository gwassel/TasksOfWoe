package usecase

import (
	"github.com/gwassel/TasksOfWoe/internal/domain/encoder"
	"github.com/gwassel/TasksOfWoe/internal/infra"
)

type Usecase struct {
	logger   infra.Logger
	taskRepo TaskRepo
	key      string
}

func New(logger infra.Logger, taskRepo TaskRepo, key string) *Usecase {
	return &Usecase{logger: logger, taskRepo: taskRepo, key: key}
}

func (u *Usecase) Handle() error {
	tasks, err := u.taskRepo.GetUnencryptedTasks()
	if err != nil {
		return err
	}

	e, err := encoder.New(u.key)
	if err != nil {
		return err
	}
	for _, task := range tasks {
		println(task.Task)
		ciphertext, err := e.Encode(task.Task)
		if err != nil {
			return err
		}
		err = u.taskRepo.EncryptTask(task.ID, ciphertext)
		if err != nil {
			return err
		}
	}

	return nil
}
