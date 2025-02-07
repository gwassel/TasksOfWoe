package complete

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gwassel/TasksOfWoe/internal/infra"
	"github.com/pkg/errors"
	"strconv"
	"strings"
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
	text := message.Text
	userID := message.From.ID

	taskID := strings.TrimSpace(strings.TrimPrefix(text, "com"))
	if taskID == "" {
		h.sendMessage(message.Chat.ID, "Please provide a task ID.")
		return
	}
	taskIDInt, err := strconv.ParseInt(taskID, 10, 64)
	if err != nil {
		h.sendMessage(message.Chat.ID, "TaskID is not int")
		return
	}
	resp := h.usecase.Handle(userID, taskIDInt)
	h.sendMessage(message.Chat.ID, resp)
}
