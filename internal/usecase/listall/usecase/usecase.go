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
	tasks, err := u.taskRepo.ListAllTasks(userID)
	if err != nil {
		return "Failed to listall tasks"
	}
	if len(tasks) == 0 {
		return "You have no tasks"
	}
	var taskList strings.Builder
	for _, task := range tasks {
		status := "Incomplete"
		if task.Completed {
			status = "Completed"
		}
		if strings.Count(task.Task, "\n") > 0 {
			taskList.WriteString(fmt.Sprintf("%d. \"%s\" [%s]\n", task.ID, task.Task, status))
		} else {
			taskList.WriteString(fmt.Sprintf("%d. %s [%s]\n", task.ID, task.Task, status))
		}
	}
	return taskList.String()
}
