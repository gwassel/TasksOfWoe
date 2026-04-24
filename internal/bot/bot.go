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
	middlewares []Middleware
}

type Middleware interface {
	Handle(message *tgbotapi.Message, next func(message *tgbotapi.Message))
}

type HandlerFunc func(message *tgbotapi.Message)

func NewBot(
	api *tgbotapi.BotAPI,
	logger infra.Logger,
	handlersMap map[string]interface {
		Handle(message *tgbotapi.Message)
	},
) *Bot {
	return &Bot{API: api, logger: logger, handlersMap: handlersMap, middlewares: []Middleware{}}
}

func (b *Bot) AddMiddleware(middleware Middleware) {
	b.middlewares = append(b.middlewares, middleware)
}

func (b *Bot) HandleMessage(message *tgbotapi.Message) {
	text := message.Text
	if text == "" {
		return
	}

	var handlerFunc HandlerFunc

	switch {
	case strings.HasPrefix(text, "add"):
		handlerFunc = b.handlersMap["add"].Handle
	case text == "ls":
		fallthrough
	case text == "list":
		fmt.Println(text)
		handlerFunc = b.handlersMap["list"].Handle
	case text == "la":
		fallthrough
	case text == "listall":
		handlerFunc = b.handlersMap["listall"].Handle
	case strings.HasPrefix(text, "com"):
		handlerFunc = b.handlersMap["com"].Handle
	case strings.HasPrefix(text, "take"):
		handlerFunc = b.handlersMap["take"].Handle
	case strings.HasPrefix(text, "untake"):
		handlerFunc = b.handlersMap["untake"].Handle
	case strings.HasPrefix(text, "desc"):
		handlerFunc = b.handlersMap["desc"].Handle
	case strings.HasPrefix(text, "help"):
		handlerFunc = b.handlersMap["help"].Handle
	case strings.HasPrefix(text, "/"):
		handlerFunc = b.handlersMap["admin"].Handle
	default:
		b.SendMessage(
			message.Chat.ID,
			`Unknown command\. Type \"help\" to see the list of available commands\.`,
		)
		return
	}

	if handlerFunc != nil {
		b.executeMiddlewareChain(message, handlerFunc)
	}
}

func (b *Bot) executeMiddlewareChain(message *tgbotapi.Message, handlerFunc HandlerFunc) {
	if len(b.middlewares) == 0 {
		handlerFunc(message)
		return
	}

	currentIndex := 0
	var next func(msg *tgbotapi.Message)

	next = func(msg *tgbotapi.Message) {
		if currentIndex < len(b.middlewares) {
			middleware := b.middlewares[currentIndex]
			currentIndex++
			middleware.Handle(msg, next)
		} else {
			handlerFunc(msg)
		}
	}

	next(message)
}

func (b *Bot) SendMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = tgbotapi.ModeMarkdownV2
	_, err := b.API.Send(msg)
	if err != nil {
		b.logger.Error(errors.Wrap(err, "unable to send message").Error())
	}
}
