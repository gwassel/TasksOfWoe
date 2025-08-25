package help

import (
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gwassel/TasksOfWoe/internal/infra"
	"github.com/pkg/errors"
)

type Handler struct {
	logger infra.Logger
	api    BotApi
}

func New(logger infra.Logger, api *tgbotapi.BotAPI) *Handler {
	return &Handler{logger: logger, api: api}
}

func (h *Handler) sendMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = tgbotapi.ModeMarkdownV2
	_, err := h.api.Send(msg)
	if err != nil {
		h.logger.Error(errors.Wrap(err, "unable to send message"))
	}
}

const messagetext string = `Available commands:
*help* \- list available commands
*add* \- add new task
*com* \(_complete_\) \- complete a task
*desc* \(_description_\) \- print task description
*ls* \(_list_\) \- list current tasks
*la* \(_listall_\) \- list all tasks
*take* \- start working on an incomplete task
*untake* \- stop working on an active task
`

func (h *Handler) Handle(message *tgbotapi.Message) {
	text := strings.TrimSpace(strings.TrimPrefix(message.Text, "help"))

	if text == "" {
		h.sendMessage(message.Chat.ID, messagetext)
	}
}
