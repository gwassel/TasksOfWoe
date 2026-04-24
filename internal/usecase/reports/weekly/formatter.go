package weekly

import (
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gwassel/TasksOfWoe/internal/domain/performance"
)

func FormatWeeklyReport(report performance.PerformanceReport) string {
	var builder strings.Builder

	endWeek := report.Date.AddDate(0, 0, 6)
	builder.WriteString(fmt.Sprintf("📅 *Weekly Performance Report*\n"))
	builder.WriteString(
		fmt.Sprintf(
			"_%s to %s_\n\n",
			report.Date.Format("2006-01-02"),
			endWeek.Format("2006-01-02"),
		),
	)

	// Overall system stats
	builder.WriteString("*🎯 System Overview*\n")
	builder.WriteString(fmt.Sprintf("Total Requests: `%d`\n", report.TotalRequests))
	if report.TotalRequests > 0 {
		builder.WriteString(
			fmt.Sprintf("Average Response Time: `%.2fms`\n\n", report.AverageSystemTime),
		)

		// Calculate daily average
		dailyAvg := float64(report.TotalRequests) / 7
		builder.WriteString(fmt.Sprintf("Avg Requests/Day: `%.1f`\n\n", dailyAvg))
	}

	// Handler statistics sorted by total requests
	handlerCount := len(report.HandlerStats)
	builder.WriteString(fmt.Sprintf("*📈 Handler Statistics* (%d handlers)\n", handlerCount))

	sortedHandlers := make([]string, 0, len(report.HandlerStats))
	for handlerName := range report.HandlerStats {
		sortedHandlers = append(sortedHandlers, handlerName)
	}

	for _, handlerName := range sortedHandlers {
		stats := report.HandlerStats[handlerName]
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

	// Anomalies detected
	anomalousHandlers := make([]string, 0)
	for handlerName, trend := range report.Trends {
		if trend.Change > 20 || trend.Change < -20 {
			anomalousHandlers = append(
				anomalousHandlers,
				fmt.Sprintf("%s: %.1f%%", handlerName, trend.Change),
			)
		}
	}

	if len(anomalousHandlers) > 0 {
		builder.WriteString("*⚠️ Anomalies Detected*\n")
		for _, anom := range anomalousHandlers {
			builder.WriteString(fmt.Sprintf("• %s\n", anom))
		}
		builder.WriteString("\n")
	}

	// Performance percentiles (top 5 handlers by traffic)
	top5Handlers := getTopHandlersByTraffic(report.HandlerStats, 5)
	builder.WriteString("*⚡ Performance Percentiles (Top 5 Handlers)*\n")

	for _, handlerName := range top5Handlers {
		if p, ok := report.Percentiles[handlerName]; ok {
			builder.WriteString(fmt.Sprintf("*%s*:\n", handlerName))
			builder.WriteString(
				fmt.Sprintf("```\nP50: %.1fms\nP90: %.1fms\nP95: %.1fms\nP99: %.1fms```\n\n",
					p.P50, p.P90, p.P95, p.P99),
			)
		}
	}

	// Slowest requests of the week
	if len(report.SlowestRequests) > 0 {
		builder.WriteString("*🐌 Slowest Requests (Top 10)*\n")
		for i := 0; i < min(len(report.SlowestRequests), 10); i++ {
			req := report.SlowestRequests[i]
			builder.WriteString(fmt.Sprintf("%d. `%s` - `%dms` (User: %d, %s)\n",
				i+1, req.Command, req.DurationMs, req.UserID, req.Timestamp.Format("15:04")))
		}
		builder.WriteString("\n")
	}

	return builder.String()
}

func getTopHandlersByTraffic(handlerStats map[string]performance.HandlerStats, limit int) []string {
	type handlerWithTraffic struct {
		name  string
		count int64
	}

	traffic := make([]handlerWithTraffic, 0, len(handlerStats))
	for name, stats := range handlerStats {
		traffic = append(traffic, handlerWithTraffic{name: name, count: stats.TotalRequests})
	}

	// Sort by traffic count
	for i := 0; i < len(traffic); i++ {
		for j := i + 1; j < len(traffic); j++ {
			if traffic[i].count < traffic[j].count {
				traffic[i], traffic[j] = traffic[j], traffic[i]
			}
		}
	}

	results := make([]string, 0, min(limit, len(traffic)))
	for i := 0; i < min(limit, len(traffic)); i++ {
		results = append(results, traffic[i].name)
	}

	return results
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func SendWeeklyReportToAdmins(
	botAPI *tgbotapi.BotAPI,
	adminUsers []int64,
	report performance.PerformanceReport,
) error {
	reportText := FormatWeeklyReport(report)

	for _, adminID := range adminUsers {
		msg := tgbotapi.NewMessage(adminID, reportText)
		msg.ParseMode = tgbotapi.ModeMarkdownV2

		_, err := botAPI.Send(msg)
		if err != nil {
			return fmt.Errorf("failed to send weekly report to admin %d: %w", adminID, err)
		}
	}

	return nil
}
