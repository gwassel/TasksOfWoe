package help

import (
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	domain "github.com/gwassel/TasksOfWoe/internal/domain/task"
	"github.com/gwassel/TasksOfWoe/internal/infra"
	"github.com/pkg/errors"
)

type Handler struct {
	logger     infra.Logger
	api        BotApi
	descsmap   map[string]domain.Description
	descsslice []domain.Description
}

func New(
	logger infra.Logger,
	api *tgbotapi.BotAPI,
	descsmap map[string]domain.Description,
	descsslice []domain.Description,
) *Handler {
	return &Handler{logger: logger, api: api, descsmap: descsmap, descsslice: descsslice}
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
	if text == "" {
		helpMessage.WriteString("Available commands:\n")
		for _, desc := range h.descsslice {
			helpMessage.WriteString(printShort(desc))
		}
		helpMessage.WriteString(`type \"help \<command\>\" for command description`)
	} else {
		desc := h.descsmap[text]
		helpMessage.WriteString(printFull(desc))
	}

	h.sendMessage(message.Chat.ID, helpMessage.String())
}

func printShort(d domain.Description) string {
	var text strings.Builder

	text.WriteString(
		fmt.Sprintf("*%s*", tgbotapi.EscapeText(tgbotapi.ModeMarkdownV2, d.Name)),
	)
	if d.Aliases != nil {
		text.WriteString(` \(_`)
		for _, alias := range d.Aliases {
			text.WriteString(alias)
		}
		text.WriteString(`_\) `)
	}
	text.WriteString(
		` \- ` + tgbotapi.EscapeText(tgbotapi.ModeMarkdownV2, d.DescShort) + "\n",
	)

	return text.String()
}

func printFull(d domain.Description) string {
	if d.Name == "" {
		return `Unknown command\. Type \"help\" to see the list of available commands\.`
	}

	var text strings.Builder

	text.WriteString(
		fmt.Sprintf("*%s*", tgbotapi.EscapeText(tgbotapi.ModeMarkdownV2, d.Name)),
	)
	if d.Aliases != nil {
		text.WriteString(` \(_`)
		for _, alias := range d.Aliases {
			text.WriteString(alias)
		}
		text.WriteString(`_\) `)
	}
	text.WriteString(
		` \- ` + tgbotapi.EscapeText(tgbotapi.ModeMarkdownV2, d.DescFull) + "\n",
	)

	text.WriteString(
		fmt.Sprintf("usage: `%s`\n", tgbotapi.EscapeText(tgbotapi.ModeMarkdownV2, d.Format)),
	)
	for _, arg := range d.Args {
		text.WriteString(tgbotapi.EscapeText(tgbotapi.ModeMarkdownV2, arg))
	}

	return text.String()
}
