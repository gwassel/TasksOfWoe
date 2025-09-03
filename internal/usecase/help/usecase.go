package usecase

import domain "github.com/gwassel/TasksOfWoe/internal/domain/task"

type Usecase struct {
	Desc domain.Description
}

func NewUsecase() *Usecase {
	desc := domain.Description{
		Name:      "help",
		Aliases:   nil,
		DescShort: "print help message",
		DescFull:  "print help message for all commands or a specific command",
		Format:    "help <command>",
		Args:      []string{"command: command to get help with (optional)"},
	}
	return &Usecase{Desc: desc}
}
