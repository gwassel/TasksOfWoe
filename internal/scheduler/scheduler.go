package scheduler

import (
	"context"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	dailyReport "github.com/gwassel/TasksOfWoe/internal/usecase/reports/daily"
	weeklyReport "github.com/gwassel/TasksOfWoe/internal/usecase/reports/weekly"
	"github.com/pkg/errors"
)

type Scheduler struct {
	dailyUsecase  *dailyReport.Usecase
	weeklyUsecase *weeklyReport.Usecase
	botAPI        *tgbotapi.BotAPI
	config        SchedulerConfig
	logger        Logger
}

type Logger interface {
	Info(args ...interface{})
	Error(args ...interface{})
	Warn(args ...interface{})
}

func NewScheduler(
	dailyUsecase *dailyReport.Usecase,
	weeklyUsecase *weeklyReport.Usecase,
	botAPI *tgbotapi.BotAPI,
	config SchedulerConfig,
	logger Logger,
) *Scheduler {
	if config.DailyReportHour < 0 || config.DailyReportHour > 23 {
		config.DailyReportHour = 9
	}
	if config.WeeklyReportHour < 0 || config.WeeklyReportHour > 23 {
		config.WeeklyReportHour = 10
	}

	return &Scheduler{
		dailyUsecase:  dailyUsecase,
		weeklyUsecase: weeklyUsecase,
		botAPI:        botAPI,
		config:        config,
		logger:        logger,
	}
}

func (s *Scheduler) Start(ctx context.Context) {
	go s.runDailyScheduler(ctx)
	go s.runWeeklyScheduler(ctx)
}

func (s *Scheduler) runDailyScheduler(ctx context.Context) {
	s.logger.Info("Daily scheduler started")

	for {
		now := time.Now()
		nextRun := s.getNextDailyRunTime(now)

		s.logger.Info("Next daily report scheduled", "time", nextRun.String())

		waitDuration := nextRun.Sub(now)

		select {
		case <-time.After(waitDuration):
			s.sendDailyReports(ctx)
		case <-ctx.Done():
			s.logger.Info("Daily scheduler stopped")
			return
		}
	}
}

func (s *Scheduler) runWeeklyScheduler(ctx context.Context) {
	s.logger.Info("Weekly scheduler started")

	for {
		now := time.Now()
		nextRun := s.getNextWeeklyRunTime(now)

		s.logger.Info("Next weekly report scheduled", "time", nextRun.String())

		waitDuration := nextRun.Sub(now)

		select {
		case <-time.After(waitDuration):
			s.sendWeeklyReports(ctx)
		case <-ctx.Done():
			s.logger.Info("Weekly scheduler stopped")
			return
		}
	}
}

func (s *Scheduler) getNextDailyRunTime(now time.Time) time.Time {
	nextRun := time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		s.config.DailyReportHour,
		0,
		0,
		0,
		now.Location(),
	)

	// If the scheduled time has already passed today, schedule for tomorrow
	if nextRun.Before(now) {
		nextRun = nextRun.Add(24 * time.Hour)
	}

	return nextRun
}

func (s *Scheduler) getNextWeeklyRunTime(now time.Time) time.Time {
	// Find the next occurrence of the configured day
	daysUntilSunday := int((time.Sunday - now.Weekday() + 7) % 7)

	nextRun := time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		s.config.WeeklyReportHour,
		0,
		0,
		0,
		now.Location(),
	)
	nextRun = nextRun.AddDate(0, 0, daysUntilSunday)

	// If the scheduled time has already passed, go to next week
	if nextRun.Before(now) {
		nextRun = nextRun.AddDate(0, 0, 7)
	}

	return nextRun
}

func (s *Scheduler) sendDailyReports(ctx context.Context) {
	s.logger.Info("Sending daily performance reports")

	reportDate := time.Now().AddDate(0, 0, -1) // Yesterday's report

	report, err := s.dailyUsecase.GenerateReport(ctx, reportDate)
	if err != nil {
		s.logger.Error(errors.Wrap(err, "generating daily report").Error())
		return
	}

	adminUsers, err := s.dailyUsecase.GetAdminUsers(ctx)
	if err != nil {
		s.logger.Error(errors.Wrap(err, "getting admin users").Error())
		return
	}

	if len(adminUsers) == 0 {
		s.logger.Warn("No admin users found, skipping report delivery")
		return
	}

	err = dailyReport.SendReportToAdmins(s.botAPI, adminUsers, report)
	if err != nil {
		s.logger.Error(errors.Wrap(err, "sending daily reports").Error())
		return
	}

	s.logger.Info("Daily reports sent successfully", "admins", len(adminUsers))
}

func (s *Scheduler) sendWeeklyReports(ctx context.Context) {
	s.logger.Info("Sending weekly performance reports")

	reportDate := time.Now() // Current week's report

	report, err := s.weeklyUsecase.GenerateReport(ctx, reportDate)
	if err != nil {
		s.logger.Error(errors.Wrap(err, "generating weekly report").Error())
		return
	}

	adminUsers, err := s.weeklyUsecase.GetAdminUsers(ctx)
	if err != nil {
		s.logger.Error(errors.Wrap(err, "getting admin users").Error())
		return
	}

	if len(adminUsers) == 0 {
		s.logger.Warn("No admin users found, skipping weekly report delivery")
		return
	}

	err = weeklyReport.SendWeeklyReportToAdmins(s.botAPI, adminUsers, report)
	if err != nil {
		s.logger.Error(errors.Wrap(err, "sending weekly reports").Error())
		return
	}

	s.logger.Info("Weekly reports sent successfully", "admins", len(adminUsers))
}

func (s *Scheduler) SendTestDailyReport(ctx context.Context) error {
	reportDate := time.Now().AddDate(0, 0, -1)

	report, err := s.dailyUsecase.GenerateReport(ctx, reportDate)
	if err != nil {
		return errors.Wrap(err, "generating test daily report")
	}

	adminUsers, err := s.dailyUsecase.GetAdminUsers(ctx)
	if err != nil {
		return errors.Wrap(err, "getting admin users for test report")
	}

	if len(adminUsers) == 0 {
		return errors.New("no admin users found")
	}

	return dailyReport.SendReportToAdmins(s.botAPI, adminUsers, report)
}

func (s *Scheduler) SendTestWeeklyReport(ctx context.Context) error {
	reportDate := time.Now()

	report, err := s.weeklyUsecase.GenerateReport(ctx, reportDate)
	if err != nil {
		return errors.Wrap(err, "generating test weekly report")
	}

	adminUsers, err := s.weeklyUsecase.GetAdminUsers(ctx)
	if err != nil {
		return errors.Wrap(err, "getting admin users for test report")
	}

	if len(adminUsers) == 0 {
		return errors.New("no admin users found")
	}

	return weeklyReport.SendWeeklyReportToAdmins(s.botAPI, adminUsers, report)
}
