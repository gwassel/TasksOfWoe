package complete

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

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

const maxlen int = 50

func cutText(s string, len int) string {
	if utf8.RuneCountInString(s) <= len {
		return s
	}
	for utf8.RuneCountInString(s) > len {
		s = s[0:strings.LastIndexFunc(s, unicode.IsSpace)]
	}
	return s + " ..."
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
		h.sendMessage(message.Chat.ID, "You have no tasks, add one")
		return
	}

	var taskList strings.Builder
	var (
		separatorflag1 = true
		separatorflag2 = true
	)

	for _, task := range tasks {
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

		task.Task = cutText(task.Task, maxlen)

		if strings.Contains(task.Task, "\n") {
			taskList.WriteString(fmt.Sprintf("\"%s\" [%s]\n", task.Task, status))
		} else {
			taskList.WriteString(fmt.Sprintf("%s [%s]\n", task.Task, status))
		}
	}
	h.sendMessage(message.Chat.ID, taskList.String())
}
