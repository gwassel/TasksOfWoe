package help

import (
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gwassel/TasksOfWoe/internal/domain"
	"github.com/gwassel/TasksOfWoe/internal/infra"
	"github.com/pkg/errors"
)

type Handler struct {
	logger infra.Logger
	api    BotApi
	descs  map[string]domain.Description
}

func New(logger infra.Logger, api *tgbotapi.BotAPI, descs map[string]domain.Description) *Handler {
	return &Handler{logger: logger, api: api, descs: descs}
}

func (h *Handler) sendMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = tgbotapi.ModeMarkdownV2
	_, err := h.api.Send(msg)
	if err != nil {
		h.logger.Error(errors.Wrap(err, "unable to send message"))
	}
}

func (h *Handler) Handle(message *tgbotapi.Message) {
	text := strings.TrimSpace(strings.TrimPrefix(message.Text, "help"))

	var helpMessage strings.Builder
	helpMessage.WriteString("Available commands:\n")
	if text == "" {
		for _, desc := range h.descs {
			helpMessage.WriteString(
				fmt.Sprintf("*%s*", tgbotapi.EscapeText(tgbotapi.ModeMarkdownV2, desc.Name)),
			)
			if desc.Aliases != nil {
				helpMessage.WriteString(` \(_`)
				for _, alias := range desc.Aliases {
					helpMessage.WriteString(alias)
				}
				helpMessage.WriteString(`\)_ `)
			}
			helpMessage.WriteString(
				` \- ` + tgbotapi.EscapeText(tgbotapi.ModeMarkdownV2, desc.DescShort),
			)
			helpMessage.WriteString("\n")
		}
	}
	h.sendMessage(message.Chat.ID, helpMessage.String())
}
