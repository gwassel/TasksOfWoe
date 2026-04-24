package middleware

import (
	"fmt"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	domain_analytics "github.com/gwassel/TasksOfWoe/internal/domain/analytics"
	"github.com/gwassel/TasksOfWoe/internal/domain/performance"
	"github.com/gwassel/TasksOfWoe/internal/infra"
)

type MetricsDaemon interface {
	Write(metric performance.HandlerMetric)
}

type Middleware interface {
	Handle(message *tgbotapi.Message, next func(message *tgbotapi.Message))
}

type MetricsMiddleware struct {
	daemon MetricsDaemon
	logger infra.Logger
	an     AnalyticsClient
}

type AnalyticsClient interface {
	Write(message domain_analytics.Event)
}

func NewMetricsMiddleware(
	daemon MetricsDaemon,
	logger infra.Logger,
	an AnalyticsClient,
) *MetricsMiddleware {
	return &MetricsMiddleware{
		daemon: daemon,
		logger: logger,
		an:     an,
	}
}

func (m *MetricsMiddleware) Handle(
	message *tgbotapi.Message,
	next func(message *tgbotapi.Message),
) {
	if message == nil || message.Text == "" {
		next(message)
		return
	}

	handlerName, command := m.extractCommandInfo(message.Text)

	startTime := time.Now()

	next(message)

	duration := time.Since(startTime)
	durationMs := int64(duration.Milliseconds())

	metric := performance.HandlerMetric{
		UserID:      message.From.ID,
		HandlerName: handlerName,
		Command:     command,
		DurationMs:  durationMs,
		Timestamp:   time.Now(),
	}

	m.daemon.Write(metric)

	if durationMs > 1000 {
		m.logger.Warn(fmt.Sprintf("Slow request: %s took %dms", command, durationMs))
	}
}

func (m *MetricsMiddleware) extractCommandInfo(text string) (string, string) {
	text = strings.TrimSpace(text)

	handlerName := "unknown"
	command := text

	parts := strings.Fields(text)
	if len(parts) > 0 {
		firstWord := strings.ToLower(parts[0])

		switch firstWord {
		case "add":
			handlerName = "add"
		case "list", "ls":
			handlerName = "list"
		case "listall", "la":
			handlerName = "listall"
		case "complete", "com":
			handlerName = "complete"
		case "desc", "description":
			handlerName = "description"
		case "help":
			handlerName = "help"
		case "take":
			handlerName = "take"
		case "untake":
			handlerName = "untake"
		case "/metrics", "/stats", "/slowest", "/testdaily", "/testweekly":
			handlerName = "admin"
		default:
			handlerName = "unknown"
		}

		if len(parts) > 0 {
			command = strings.Join(parts, " ")
		}
	}

	return handlerName, command
}
