package list

import (
	"fmt"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gwassel/TasksOfWoe/internal/domain/analytics"
	"github.com/gwassel/TasksOfWoe/internal/infra"
	"github.com/pkg/errors"
)

type Handler struct {
	logger    infra.Logger
	an        AnalyticsClient
	api       BotApi
	usecase   Usecase
	maxlen    int
	mincutlen int
}

func New(logger infra.Logger, an AnalyticsClient, api *tgbotapi.BotAPI, usecase Usecase) *Handler {
	return &Handler{logger: logger, an: an, api: api, usecase: usecase, maxlen: 50, mincutlen: 30}
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

	h.an.Write(analytics.NewEvent(userID, "list tasks", time.Now()))

	var taskList strings.Builder
	separatorFlag := true

	for _, task := range tasks {
		if task.InWork {
			fmt.Fprintf(&taskList, "%d*. ", task.UserTaskID)
		} else {
			if separatorFlag {
				separatorFlag = false
				taskList.WriteString("\n")
			}
			fmt.Fprintf(&taskList, "%d. ", task.UserTaskID)
		}

		if utf8.RuneCountInString(task.Task) > h.maxlen {
			task.Task = h.cutText(task.Task)
			// TODO: добавить кликабельность (#41)
		}

		if strings.Contains(task.Task, "\n") {
			fmt.Fprintf(&taskList, "\"%s\"\n", task.Task)
		} else {
			fmt.Fprintf(&taskList, "%s\n", task.Task)
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
