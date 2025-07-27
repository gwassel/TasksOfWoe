package usecase

import (
	"fmt"
	"github.com/gwassel/TasksOfWoe/internal/infra"
	"github.com/pkg/errors"
	"strings"
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
		// TODO: group output by status
		taskList.WriteString(fmt.Sprintf("%d", task.UserTaskID))
		if task.InWork {
			taskList.WriteString("*")
		}
		taskList.WriteString(". ")
		if strings.Contains(task.Task, "\n") {
			taskList.WriteString(fmt.Sprintf("\"%s\"\n", task.Task))
		} else {
			taskList.WriteString(fmt.Sprintf("%s\n", task.Task))
		}
	}
	return taskList.String()
}
