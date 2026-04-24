package performance

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/gwassel/TasksOfWoe/internal/domain/performance"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type MetricsRepository interface {
	StoreMetric(ctx context.Context, metric performance.HandlerMetric) error
	StoreMetricsBatch(ctx context.Context, metrics []performance.HandlerMetric) error
	GetHandlerStats(
		ctx context.Context,
		handlerName string,
		startTime, endTime time.Time,
	) (performance.HandlerStats, error)
	GetSlowestRequests(
		ctx context.Context,
		limit int,
		startTime, endTime time.Time,
	) ([]performance.HandlerMetric, error)
	GetPercentiles(
		ctx context.Context,
		handlerName string,
		startTime, endTime time.Time,
	) (performance.PercentileStats, error)
	GetAllHandlerStats(
		ctx context.Context,
		startTime, endTime time.Time,
	) (map[string]performance.HandlerStats, error)
	DeleteOldMetrics(ctx context.Context, olderThan time.Time) (int64, error)
}

type repository struct {
	db *sqlx.DB
}

func NewMetricsRepository(db *sqlx.DB) MetricsRepository {
	return &repository{db: db}
}

func (r *repository) StoreMetric(ctx context.Context, metric performance.HandlerMetric) error {
	query := `INSERT INTO handler_performance_metrics (user_id, handler_name, command, duration_ms, timestamp)
	          VALUES ($1, $2, $3, $4, $5)`

	_, err := r.db.ExecContext(ctx, query,
		metric.UserID, metric.HandlerName, metric.Command,
		metric.DurationMs, metric.Timestamp)
	if err != nil {
		return errors.Wrap(err, "storing performance metric")
	}

	return nil
}

func (r *repository) StoreMetricsBatch(
	ctx context.Context,
	metrics []performance.HandlerMetric,
) error {
	if len(metrics) == 0 {
		return nil
	}

	queryBuilder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Insert("handler_performance_metrics").
		Columns("user_id", "handler_name", "command", "duration_ms", "timestamp")

	for _, m := range metrics {
		queryBuilder = queryBuilder.Values(
			m.UserID, m.HandlerName, m.Command, m.DurationMs, m.Timestamp)
	}

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return errors.Wrap(err, "building batch insert query")
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return errors.Wrap(err, "executing batch insert")
	}

	return nil
}

func (r *repository) GetHandlerStats(
	ctx context.Context,
	handlerName string,
	startTime, endTime time.Time,
) (performance.HandlerStats, error) {
	query := `
		SELECT 
			COUNT(*) as total_requests,
			AVG(duration_ms) as avg_duration,
			MIN(duration_ms) as min_duration,
			MAX(duration_ms) as max_duration,
			SUM(duration_ms) as total_duration
		FROM handler_performance_metrics
		WHERE handler_name = $1 
			AND timestamp >= $2 
			AND timestamp <= $3
	`

	var stats performance.HandlerStats
	err := r.db.QueryRowContext(ctx, query, handlerName, startTime, endTime).Scan(
		&stats.TotalRequests,
		&stats.AverageDuration,
		&stats.MinDuration,
		&stats.MaxDuration,
		&stats.TotalDuration,
	)
	if err != nil {
		return stats, errors.Wrap(err, "getting handler stats")
	}

	return stats, nil
}

func (r *repository) GetSlowestRequests(
	ctx context.Context,
	limit int,
	startTime, endTime time.Time,
) ([]performance.HandlerMetric, error) {
	query := `
		SELECT user_id, handler_name, command, duration_ms, timestamp
		FROM handler_performance_metrics
		WHERE timestamp >= $1 AND timestamp <= $2
		ORDER BY duration_ms DESC
		LIMIT $3
 `

	var metrics []performance.HandlerMetric
	err := r.db.SelectContext(ctx, &metrics, query, startTime, endTime, limit)
	if err != nil {
		return nil, errors.Wrap(err, "getting slowest requests")
	}

	return metrics, nil
}

func (r *repository) GetPercentiles(
	ctx context.Context,
	handlerName string,
	startTime, endTime time.Time,
) (performance.PercentileStats, error) {
	// Get all durations for the handler
	query := `
		SELECT duration_ms
		FROM handler_performance_metrics
		WHERE handler_name = $1 
			AND timestamp >= $2 
			AND timestamp <= $3
		ORDER BY duration_ms
 	`

	var durations []int64
	err := r.db.SelectContext(ctx, &durations, query, handlerName, startTime, endTime)
	if err != nil {
		return performance.PercentileStats{}, errors.Wrap(err, "getting durations for percentiles")
	}

	if len(durations) == 0 {
		return performance.PercentileStats{}, nil
	}

	// Calculate percentiles
	return calculatePercentiles(durations), nil
}

func (r *repository) GetAllHandlerStats(
	ctx context.Context,
	startTime, endTime time.Time,
) (map[string]performance.HandlerStats, error) {
	query := `
		SELECT 
			handler_name,
			COUNT(*) as total_requests,
			AVG(duration_ms) as avg_duration,
			MIN(duration_ms) as min_duration,
			MAX(duration_ms) as max_duration,
			SUM(duration_ms) as total_duration
		FROM handler_performance_metrics
		WHERE timestamp >= $1 AND timestamp <= $2
		GROUP BY handler_name
		ORDER BY total_requests DESC
 `

	type handlerStatsRow struct {
		HandlerName     string  `db:"handler_name"`
		TotalRequests   int64   `db:"total_requests"`
		AverageDuration float64 `db:"avg_duration"`
		MinDuration     int64   `db:"min_duration"`
		MaxDuration     int64   `db:"max_duration"`
		TotalDuration   int64   `db:"total_duration"`
	}

	var rows []handlerStatsRow
	err := r.db.SelectContext(ctx, &rows, query, startTime, endTime)
	if err != nil {
		return nil, errors.Wrap(err, "getting all handler stats")
	}

	stats := make(map[string]performance.HandlerStats)
	for _, row := range rows {
		hstats := performance.HandlerStats{
			TotalRequests:   row.TotalRequests,
			AverageDuration: row.AverageDuration,
			MinDuration:     row.MinDuration,
			MaxDuration:     row.MaxDuration,
			TotalDuration:   row.TotalDuration,
		}
		stats[row.HandlerName] = hstats
	}

	return stats, nil
}

func (r *repository) DeleteOldMetrics(ctx context.Context, olderThan time.Time) (int64, error) {
	query := `DELETE FROM handler_performance_metrics WHERE timestamp < $1`

	result, err := r.db.ExecContext(ctx, query, olderThan)
	if err != nil {
		return 0, errors.Wrap(err, "deleting old metrics")
	}

	return result.RowsAffected()
}

func calculatePercentiles(durations []int64) performance.PercentileStats {
	n := len(durations)
	if n == 0 {
		return performance.PercentileStats{}
	}

	getPercentile := func(p float64) int64 {
		index := int(p * float64(n-1))
		if index < 0 {
			index = 0
		}
		if index >= n {
			index = n - 1
		}
		return durations[index]
	}

	return performance.PercentileStats{
		P50: float64(getPercentile(0.50)),
		P75: float64(getPercentile(0.75)),
		P90: float64(getPercentile(0.90)),
		P95: float64(getPercentile(0.95)),
		P99: float64(getPercentile(0.99)),
	}
}
