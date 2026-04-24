package daily

import (
	"context"
	"fmt"
	"time"

	"github.com/gwassel/TasksOfWoe/internal/domain/performance"
	perfpersistence "github.com/gwassel/TasksOfWoe/internal/persistence/performance"
	"github.com/pkg/errors"
)

type Usecase struct {
	metricsRepo perfpersistence.MetricsRepository
	adminRepo   perfpersistence.AdminRepository
	logger      Logger
}

type Logger interface {
	Info(args ...interface{})
	Error(args ...interface{})
	Warn(args ...interface{})
}

func NewUsecase(
	metricsRepo perfpersistence.MetricsRepository,
	adminRepo perfpersistence.AdminRepository,
	logger Logger,
) *Usecase {
	return &Usecase{
		metricsRepo: metricsRepo,
		adminRepo:   adminRepo,
		logger:      logger,
	}
}

func (u *Usecase) GenerateReport(
	ctx context.Context,
	reportDate time.Time,
) (performance.PerformanceReport, error) {
	// Calculate time range for the report
	endTime := reportDate.Add(24 * time.Hour)
	startTime := reportDate

	// Get all handler stats
	handlerStats, err := u.metricsRepo.GetAllHandlerStats(ctx, startTime, endTime)
	if err != nil {
		return performance.PerformanceReport{}, errors.Wrap(err, "getting handler stats")
	}

	// Get slowest requests
	slowestRequests, err := u.metricsRepo.GetSlowestRequests(ctx, 10, startTime, endTime)
	if err != nil {
		u.logger.Warn(errors.Wrap(err, "getting slowest requests").Error())
		slowestRequests = []performance.HandlerMetric{}
	}

	// Calculate percentiles for each handler
	percentiles := make(map[string]performance.PercentileStats)
	for handlerName := range handlerStats {
		p, err := u.metricsRepo.GetPercentiles(ctx, handlerName, startTime, endTime)
		if err != nil {
			u.logger.Warn(
				errors.Wrap(err, fmt.Sprintf("getting percentiles for %s", handlerName)).Error(),
			)
			continue
		}
		percentiles[handlerName] = p
	}

	// Calculate system-wide stats
	totalRequests := int64(0)
	totalDuration := int64(0)
	for _, stat := range handlerStats {
		totalRequests += stat.TotalRequests
		totalDuration += stat.TotalDuration
	}

	averageSystemTime := float64(0)
	if totalRequests > 0 {
		averageSystemTime = float64(totalDuration) / float64(totalRequests)
	}

	// Calculate trends (compare with previous day)
	previousDay := startTime.AddDate(0, 0, -1)
	previousDayEnd := startTime
	trends := make(map[string]performance.TrendData)

	for handlerName, currentStats := range handlerStats {
		prevStats, err := u.metricsRepo.GetHandlerStats(
			ctx,
			handlerName,
			previousDay,
			previousDayEnd,
		)
		if err != nil || prevStats.TotalRequests == 0 {
			continue
		}

		trend := performance.TrendData{
			CurrentAvg:  currentStats.AverageDuration,
			PreviousAvg: prevStats.AverageDuration,
		}

		if prevStats.AverageDuration > 0 {
			trend.Change = ((currentStats.AverageDuration - prevStats.AverageDuration) / prevStats.AverageDuration) * 100
		}

		if trend.Change > 5 {
			trend.Direction = "up"
		} else if trend.Change < -5 {
			trend.Direction = "down"
		} else {
			trend.Direction = "stable"
		}

		trends[handlerName] = trend
	}

	return performance.PerformanceReport{
		Date:              reportDate,
		Period:            "daily",
		HandlerStats:      handlerStats,
		SlowestRequests:   slowestRequests,
		Percentiles:       percentiles,
		Trends:            trends,
		TotalRequests:     totalRequests,
		AverageSystemTime: averageSystemTime,
	}, nil
}

func (u *Usecase) GetAdminUsers(ctx context.Context) ([]int64, error) {
	return u.adminRepo.GetAdminUsers(ctx)
}

func (u *Usecase) CleanupOldMetrics(ctx context.Context, retentionDays int) error {
	if retentionDays <= 0 {
		return nil
	}

	cutoffTime := time.Now().AddDate(0, 0, -retentionDays)
	rowsAffected, err := u.metricsRepo.DeleteOldMetrics(ctx, cutoffTime)
	if err != nil {
		return err
	}

	u.logger.Info(
		"Cleaned up old performance metrics",
		"rows_affected",
		rowsAffected,
		"cutoff",
		cutoffTime.String(),
	)
	return nil
}
