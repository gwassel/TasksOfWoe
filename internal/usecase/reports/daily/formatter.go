package daily

import (
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gwassel/TasksOfWoe/internal/domain/performance"
)

func FormatReport(report performance.PerformanceReport) string {
	var builder strings.Builder

	builder.WriteString(
		fmt.Sprintf("📊 *Performance Report* (%s)\n", report.Date.Format("2006-01-02")),
	)
	builder.WriteString(fmt.Sprintf("_Period: %s_\n\n", report.Period))

	// Overall system stats
	builder.WriteString("*🎯 System Overview*\n")
	builder.WriteString(fmt.Sprintf("Total Requests: `%d`\n", report.TotalRequests))
	if report.TotalRequests > 0 {
		builder.WriteString(
			fmt.Sprintf("Average Response Time: `%.2fms`\n\n", report.AverageSystemTime),
		)
	}

	// Handler statistics
	builder.WriteString("*📈 Handler Statistics*\n")
	handlerCount := 0
	for handlerName, stats := range report.HandlerStats {
		handlerCount++
		builder.WriteString(fmt.Sprintf("• *%s*: `%d` requests, `%.2fms` avg",
			handlerName, stats.TotalRequests, stats.AverageDuration))

		if trend, ok := report.Trends[handlerName]; ok {
			switch trend.Direction {
			case "up":
				builder.WriteString(fmt.Sprintf(" 🔺%.1f%%", trend.Change))
			case "down":
				builder.WriteString(fmt.Sprintf(" 🔻%.1f%%", trend.Change))
			}
		}
		builder.WriteString("\n")
	}
	builder.WriteString("\n")

	// Percentiles
	builder.WriteString("*⚡ Performance Percentiles*\n")
	for handlerName, p := range report.Percentiles {
		builder.WriteString(fmt.Sprintf("*%s*:\n", handlerName))
		builder.WriteString(
			fmt.Sprintf(
				"```\nP50: %.1fms\nP75: %.1fms\nP90: %.1fms\nP95: %.1fms\nP99: %.1fms```\n\n",
				p.P50,
				p.P75,
				p.P90,
				p.P95,
				p.P99,
			),
		)
	}

	// Slowest requests
	if len(report.SlowestRequests) > 0 {
		builder.WriteString("*🐌 Slowest Requests*\n")
		for i, req := range report.SlowestRequests {
			builder.WriteString(fmt.Sprintf("%d. `%s` - `%dms` (User: %d)\n",
				i+1, req.Command, req.DurationMs, req.UserID))
		}
		builder.WriteString("\n")
	}

	return builder.String()
}

func SendReportToAdmins(
	botAPI *tgbotapi.BotAPI,
	adminUsers []int64,
	report performance.PerformanceReport,
) error {
	reportText := FormatReport(report)

	for _, adminID := range adminUsers {
		msg := tgbotapi.NewMessage(adminID, reportText)
		msg.ParseMode = tgbotapi.ModeMarkdownV2

		_, err := botAPI.Send(msg)
		if err != nil {
			return fmt.Errorf("failed to send report to admin %d: %w", adminID, err)
		}
	}

	return nil
}
