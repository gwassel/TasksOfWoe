package complete

import (
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gwassel/TasksOfWoe/internal/infra"
	"github.com/pkg/errors"
)

type Handler struct {
	logger  infra.Logger
	api     BotApi
	usecase Usecase
}

func New(logger infra.Logger, api *tgbotapi.BotAPI, usecase Usecase) *Handler {
	return &Handler{logger: logger, api: api, usecase: usecase}
}

func (h *Handler) sendMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	_, err := h.api.Send(msg)
	if err != nil {
		h.logger.Error(errors.Wrap(err, "unable to send message"))
	}
}

func (h *Handler) Handle(message *tgbotapi.Message) {
	userID := message.From.ID

	tasks, err := h.usecase.Handle(userID)
	if err != nil {
		h.logger.Error(err)
		h.sendMessage(message.Chat.ID, "Unable to list tasks.")
		return
	}
	if len(tasks) == 0 {
		h.sendMessage(message.Chat.ID, "You have nothing to do")
		return
	}

	var taskList strings.Builder
	var separatorflag = true

	for _, task := range tasks {
		if task.InWork {
			taskList.WriteString(fmt.Sprintf("%d*. ", task.UserTaskID))
		} else {
			if separatorflag {
				separatorflag = false
				taskList.WriteString("\n")
			}
			taskList.WriteString(fmt.Sprintf("%d. ", task.UserTaskID))
		}

		if strings.Contains(task.Task, "\n") {
			taskList.WriteString(fmt.Sprintf("\"%s\"\n", task.Task))
		} else {
			taskList.WriteString(fmt.Sprintf("%s\n", task.Task))
		}
	}
	h.sendMessage(message.Chat.ID, taskList.String())
}
