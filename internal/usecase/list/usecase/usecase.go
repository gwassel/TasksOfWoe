package usecase

import (
	"fmt"
	"strings"
)

type Usecase struct {
	taskRepo TaskRepo
}

func New(taskRepo TaskRepo) *Usecase {
	return &Usecase{taskRepo: taskRepo}
}

func (u *Usecase) Handle(userID int64) string {
	tasks, err := u.taskRepo.ListTasks(userID)
	if err != nil {
		return "Failed to list tasks"
	}
	if len(tasks) == 0 {
		return "You have nothing to do"
	}

	var taskList strings.Builder
	for _, task := range tasks {
		taskList.WriteString(fmt.Sprintf("%d. %s\n", task.ID, task.Task))
	}
	return taskList.String()
}
