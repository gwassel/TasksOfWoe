package complete

import (
	"fmt"
	"strconv"
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
	text := message.Text
	userID := message.From.ID

	taskID := strings.TrimSpace(strings.TrimPrefix(text, "untake"))
	if taskID == "" {
		h.sendMessage(message.Chat.ID, "Please provide a task ID.")
		return
	}
	h.logger.Debug(fmt.Sprintf("taskID: %s", taskID))

	//taskIDInt, err := strconv.ParseInt(taskID, 10, 64)
	taskIDs, err := h.convertInput(taskID)
	if err != nil {
		h.sendMessage(message.Chat.ID, "TaskID is not int")
		return
	}
	err = h.usecase.Handle(userID, taskIDs, false)
	if err != nil {
		h.logger.Error(err)
		h.sendMessage(message.Chat.ID, "Unable to update task status.")
		return
	}
	h.sendMessage(message.Chat.ID, "Suspended work on task(s).")
}

func (h *Handler) convertInput(strIDs string) ([]int64, error) {
	result := make([]int64, 0)
	for _, k := range strings.Split(strIDs, " ") {
		tmp, err := strconv.ParseInt(k, 10, 64)
		if err != nil {
			return nil, err
		}
		result = append(result, tmp)
	}

	return result, nil
}
