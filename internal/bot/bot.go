package bot

import (
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gwassel/TasksOfWoe/internal/infra"
	"github.com/pkg/errors"
)

type Bot struct {
	API         *tgbotapi.BotAPI
	logger      infra.Logger
	handlersMap map[string]interface {
		Handle(message *tgbotapi.Message)
	}
}

func NewBot(
	api *tgbotapi.BotAPI,
	logger infra.Logger,
	handlersMap map[string]interface {
		Handle(message *tgbotapi.Message)
	},
) *Bot {
	return &Bot{API: api, logger: logger, handlersMap: handlersMap}
}

func (b *Bot) HandleMessage(message *tgbotapi.Message) {
	text := message.Text

	switch {
	case strings.HasPrefix(text, "add"):
		b.handlersMap["add"].Handle(message)

	case text == "ls":
		fallthrough
	case text == "list":
		fmt.Println(text)
		b.handlersMap["list"].Handle(message)

	case text == "la":
		fallthrough
	case text == "listall":
		b.handlersMap["listall"].Handle(message)

	case strings.HasPrefix(text, "com"):
		b.handlersMap["com"].Handle(message)

	case strings.HasPrefix(text, "take"):
		b.handlersMap["take"].Handle(message)

	case strings.HasPrefix(text, "untake"):
		b.handlersMap["untake"].Handle(message)

	case strings.HasPrefix(text, "desc"):
		b.handlersMap["desc"].Handle(message)

	case strings.HasPrefix(text, "help"):
		b.handlersMap["help"].Handle(message)

	default:
		b.SendMessage(message.Chat.ID, `Unknown command\. Type \"help\" to see the list of available commands\.`)
	}
}

func (b *Bot) SendMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = tgbotapi.ModeMarkdownV2
	_, err := b.API.Send(msg)
	if err != nil {
		b.logger.Error(errors.Wrap(err, "unable to send message"))
	}
}
