package complete

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gwassel/TasksOfWoe/internal/domain"
	"github.com/gwassel/TasksOfWoe/internal/infra"
	"github.com/pkg/errors"
)

type Handler struct {
	logger    infra.Logger
	api       BotApi
	usecase   Usecase
	maxlen    int
	mincutlen int
}

func New(logger infra.Logger, api *tgbotapi.BotAPI, usecase Usecase) *Handler {
	return &Handler{logger: logger, api: api, usecase: usecase, maxlen: 40, mincutlen: 20}
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
		h.sendMessage(message.Chat.ID, "You have no tasks, add one")
		return
	}

	var taskList strings.Builder
	var (
		separatorflag1 = true
		separatorflag2 = true
	)

	for _, task := range tasks {
		status := task.Status()
		switch status {
		case domain.Working:
			taskList.WriteString(fmt.Sprintf("%d*. ", task.ID))

		case domain.Incomplete:
			if separatorflag1 {
				separatorflag1 = false
				taskList.WriteString("\n")
			}
			taskList.WriteString(fmt.Sprintf("%d. ", task.ID))

		case domain.Completed:
			if separatorflag2 {
				separatorflag2 = false
				taskList.WriteString("\n")
			}
			taskList.WriteString(fmt.Sprintf("%d. ", task.ID))
		}

		if utf8.RuneCountInString(task.Task) > h.maxlen {
			task.Task = h.cutText(task.Task)
			// TODO:добавить кликабельность
		}

		if strings.Contains(task.Task, "\n") {
			taskList.WriteString(fmt.Sprintf("\"%s\" [%s]\n", task.Task, domain.ToString(status)))
		} else {
			taskList.WriteString(fmt.Sprintf("%s [%s]\n", task.Task, domain.ToString(status)))
		}
	}
	h.sendMessage(message.Chat.ID, taskList.String())
}

func (h *Handler) cutText(s string) string {
	if utf8.RuneCountInString(s) > h.maxlen {
		// Unicode compatibility
		s = string([]rune(s)[0:h.maxlen])
		cutpos := strings.LastIndexFunc(s, unicode.IsSpace)
		if cutpos != -1 {
			t := s[0:cutpos]
			if len([]rune(t)) >= h.mincutlen {
				s = t
			}
		}
		s += " ..."
	}

	return s
}
