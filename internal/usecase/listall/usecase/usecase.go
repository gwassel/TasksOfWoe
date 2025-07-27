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
		// TODO: group output by status
		status := "Incomplete"
		if task.Completed {
			status = "Completed"
		} else if task.InWork {
			status = "Working"
		}
		taskList.WriteString(fmt.Sprintf("%d", task.ID))
		if !task.Completed && task.InWork {
			taskList.WriteString("*")
		}
		taskList.WriteString(". ")
		if strings.Contains(task.Task, "\n") {
			taskList.WriteString(fmt.Sprintf("\"%s\" [%s]\n", task.Task, status))
		} else {
			taskList.WriteString(fmt.Sprintf("%s [%s]\n", task.Task, status))
		}
	}
	return taskList.String()
}
