package add

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
	text := message.Text

	task := strings.TrimSpace(strings.TrimPrefix(text, "add"))
	if task == "" {
		h.sendMessage(message.Chat.ID, "Please provide a task description.")
		return
	}
	userTaskID, err := h.usecase.Handle(userID, task)
	if err != nil {
		h.logger.Error(err)
		h.sendMessage(message.Chat.ID, "Failed to add task.")
		return
	}
	h.sendMessage(message.Chat.ID, fmt.Sprintf("Task %d added.", userTaskID))
}
