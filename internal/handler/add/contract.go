package add

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gwassel/TasksOfWoe/internal/domain/analytics"
)

type BotApi interface {
	Send(c tgbotapi.Chattable) (tgbotapi.Message, error)
}

type Usecase interface {
	Handle(ctx context.Context, userID int64, task string) (int64, error)
}

type AnalyticsClient interface {
	Write(message analytics.Event)
}
