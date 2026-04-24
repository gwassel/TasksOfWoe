package admin

import (
	"context"
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
)

type Handler struct {
	scheduler Scheduler
	logger    Logger
	api       *tgbotapi.BotAPI
}

type Scheduler interface {
	SendTestDailyReport(ctx context.Context) error
	SendTestWeeklyReport(ctx context.Context) error
}

type Logger interface {
	Info(args ...interface{})
	Error(args ...interface{})
}

type AnalyticsClient interface {
	Write(message interface{})
}

func New(logger Logger, an AnalyticsClient, api *tgbotapi.BotAPI, scheduler Scheduler) *Handler {
	return &Handler{
		scheduler: scheduler,
		logger:    logger,
		api:       api,
	}
}

func (h *Handler) sendMessage(chatID int64, text string) {
	if h.api == nil {
		return
	}

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = tgbotapi.ModeMarkdownV2
	_, err := h.api.Send(msg)
	if err != nil {
		h.logger.Error(errors.Wrap(err, "unable to send message").Error())
	}
}

func (h *Handler) Handle(message *tgbotapi.Message) {
	ctx := context.Background()

	if message == nil || message.Text == "" {
		return
	}

	parts := strings.Fields(message.Text)
	if len(parts) == 0 {
		return
	}

	command := strings.ToLower(parts[0])

	switch command {
	case "/metrics", "/stats":
		h.handleMetricsCommand(ctx, message)
	case "/slowest":
		h.handleSlowestCommand(ctx, message)
	case "/testdaily":
		h.handleTestDailyCommand(ctx, message)
	case "/testweekly":
		h.handleTestWeeklyCommand(ctx, message)
	case "/help":
		h.handleHelpCommand(message)
	default:
		h.sendMessage(
			message.Chat.ID,
			"Unknown admin command `\\`\\`\nType /help for available commands",
		)
	}
}

func (h *Handler) handleMetricsCommand(ctx context.Context, message *tgbotapi.Message) {
	h.sendMessage(message.Chat.ID, "Generating daily performance report\\.\\.\\.")

	err := h.scheduler.SendTestDailyReport(ctx)
	if err != nil {
		h.logger.Error(errors.Wrap(err, "sending test daily report").Error())
		h.sendMessage(
			message.Chat.ID,
			fmt.Sprintf("Error generating report: %s", escapeMarkdownV2(err.Error())),
		)
		return
	}

	h.sendMessage(message.Chat.ID, "Daily report sent successfully")
}

func (h *Handler) handleSlowestCommand(ctx context.Context, message *tgbotapi.Message) {
	h.sendMessage(
		message.Chat.ID,
		"Slowest requests are included in the daily report\\.\\nUse /metrics to see the full report",
	)
}

func (h *Handler) handleTestDailyCommand(ctx context.Context, message *tgbotapi.Message) {
	h.sendMessage(message.Chat.ID, "Sending test daily performance report\\.\\.\\.")

	err := h.scheduler.SendTestDailyReport(ctx)
	if err != nil {
		h.logger.Error(errors.Wrap(err, "sending test daily report").Error())
		h.sendMessage(message.Chat.ID, fmt.Sprintf("Error: %s", escapeMarkdownV2(err.Error())))
		return
	}

	h.sendMessage(message.Chat.ID, "Test daily report sent successfully")
}

func (h *Handler) handleTestWeeklyCommand(ctx context.Context, message *tgbotapi.Message) {
	h.sendMessage(message.Chat.ID, "Sending test weekly performance report\\.\\.\\.")

	err := h.scheduler.SendTestWeeklyReport(ctx)
	if err != nil {
		h.logger.Error(errors.Wrap(err, "sending test weekly report").Error())
		h.sendMessage(message.Chat.ID, fmt.Sprintf("Error: %s", escapeMarkdownV2(err.Error())))
		return
	}

	h.sendMessage(message.Chat.ID, "Test weekly report sent successfully")
}

func (h *Handler) handleHelpCommand(message *tgbotapi.Message) {
	helpText := "*Admin Commands*\n\n" +
		"• /metrics or /stats - Send daily performance report\n" +
		"• /slowest - Show information about slowest requests\n" +
		"• /testdaily - Send test daily report\n" +
		"• /testweekly - Send test weekly report\n" +
		"• /help - Show this help message"

	h.sendMessage(message.Chat.ID, helpText)
}

func escapeMarkdownV2(text string) string {
	replacer := strings.NewReplacer(
		"_", "\\_",
		"*", "\\*",
		"[", "\\[",
		"]", "\\]",
		"(", "\\(",
		")", "\\)",
		"~", "\\~",
		"`", "\\`",
		">", "\\>",
		"#", "\\#",
		"+", "\\+",
		"-", "\\-",
		"=", "\\=",
		"|", "\\|",
		"{", "\\{",
		"}", "\\}",
		".", "\\.",
		"!", "\\!",
	)

	return replacer.Replace(text)
}
