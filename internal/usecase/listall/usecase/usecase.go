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
	var (
		separatorflag1 = true
		separatorflag2 = true
	)

	for _, task := range tasks {
		// TEST: проверить корректность числа пустых строк
		status := "Incomplete"
		if task.Completed {
			status = "Completed"
		} else if task.InWork {
			status = "Working"
		}

		switch status {
		case "Working":
			taskList.WriteString(fmt.Sprintf("%d*. ", task.ID))

		case "Incomplete":
			if separatorflag1 {
				separatorflag1 = false
				taskList.WriteString("\n")
			}
			taskList.WriteString(fmt.Sprintf("%d. ", task.ID))

		case "Completed":
			if separatorflag2 {
				separatorflag2 = false
				taskList.WriteString("\n")
			}
			taskList.WriteString(fmt.Sprintf("%d. ", task.ID))
		}

		if strings.Contains(task.Task, "\n") {
			taskList.WriteString(fmt.Sprintf("\"%s\" [%s]\n", task.Task, status))
		} else {
			taskList.WriteString(fmt.Sprintf("%s [%s]\n", task.Task, status))
		}
	}
	return taskList.String()
}
