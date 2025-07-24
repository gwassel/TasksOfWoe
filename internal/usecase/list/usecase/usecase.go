package usecase

import (
	"fmt"
	"strings"

	"github.com/gwassel/TasksOfWoe/internal/infra"
	"github.com/pkg/errors"
)

type Usecase struct {
	logger   infra.Logger
	taskRepo TaskRepo
}

func New(logger infra.Logger, taskRepo TaskRepo) *Usecase {
	return &Usecase{logger: logger, taskRepo: taskRepo}
}

func (u *Usecase) Handle(userID int64) string {
	tasks, err := u.taskRepo.ListTasks(userID)
	if err != nil {
		u.logger.Error(errors.Wrap(err, "unable to list tasks"))
		return "Failed to list tasks"
	}
	if len(tasks) == 0 {
		return "You have nothing to do"
	}

	var taskList strings.Builder
	for _, task := range tasks {
		taskList.WriteString(fmt.Sprintf("%d. %s\n", task.UserTaskID, task.Task))
	}
	return taskList.String()
}
