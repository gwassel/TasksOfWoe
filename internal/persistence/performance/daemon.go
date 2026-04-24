package performance

import (
	"context"
	"sync"
	"time"

	"github.com/gwassel/TasksOfWoe/internal/domain/performance"
	"github.com/gwassel/TasksOfWoe/internal/infra"
	"github.com/pkg/errors"
)

type MetricsDaemon struct {
	db            MetricsRepository
	messageQueue  chan performance.HandlerMetric
	logger        infra.Logger
	batchSize     int
	flushInterval time.Duration
	mu            sync.Mutex
}

func NewMetricsDaemon(
	db MetricsRepository,
	queueSize, batchSize int,
	logger infra.Logger,
	flushInterval time.Duration,
) *MetricsDaemon {
	return &MetricsDaemon{
		db:            db,
		messageQueue:  make(chan performance.HandlerMetric, queueSize),
		logger:        logger,
		batchSize:     batchSize,
		flushInterval: flushInterval,
	}
}

func (d *MetricsDaemon) Write(metric performance.HandlerMetric) {
	go func() {
		d.messageQueue <- metric
	}()
}

func (d *MetricsDaemon) StartWorker(ctx context.Context) {
	go func() {
		batch := make([]performance.HandlerMetric, 0, d.batchSize)
		ticker := time.NewTicker(d.flushInterval)
		defer ticker.Stop()

		for {
			select {
			case metric := <-d.messageQueue:
				batch = append(batch, metric)
				if len(batch) >= d.batchSize {
					d.flushBatch(ctx, batch)
					batch = make([]performance.HandlerMetric, 0, d.batchSize)
				}

			case <-ticker.C:
				d.flushBatch(ctx, batch)
				batch = make([]performance.HandlerMetric, 0, d.batchSize)

			case <-ctx.Done():
				// Flush remaining batch before shutdown
				d.flushBatch(ctx, batch)
				return
			}
		}
	}()
}

func (d *MetricsDaemon) flushBatch(ctx context.Context, batch []performance.HandlerMetric) {
	if len(batch) == 0 {
		return
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	err := d.db.StoreMetricsBatch(ctx, batch)
	if err != nil {
		d.logger.Error(errors.Wrap(err, "failed to flush metrics batch").Error())
	}
}

func (d *MetricsDaemon) CleanupOldMetrics(ctx context.Context, retentionDays int) {
	if retentionDays <= 0 {
		return
	}

	cutoffTime := time.Now().AddDate(0, 0, -retentionDays)

	go func() {
		rowsAffected, err := d.db.DeleteOldMetrics(ctx, cutoffTime)
		if err != nil {
			d.logger.Error(errors.Wrap(err, "failed to delete old metrics").Error())
			return
		}

		if rowsAffected > 0 {
			d.logger.Info(
				"Deleted old performance metrics",
				"rows_affected",
				rowsAffected,
				"cutoff",
				cutoffTime.String(),
			)
		}
	}()
}
