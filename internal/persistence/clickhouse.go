package persistence

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/gwassel/TasksOfWoe/internal/domain/analytics"
	"github.com/gwassel/TasksOfWoe/internal/infra"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type CH struct {
	db           *sqlx.DB
	messageQueue chan analytics.Event
	logger       infra.Logger
}

func NewAnalyticsDaemon(db *sqlx.DB, queueSize int, logger infra.Logger) *CH {
	return &CH{
		db:           db,
		messageQueue: make(chan analytics.Event, queueSize),
		logger:       logger,
	}
}

func (c *CH) WriteToDB(ctx context.Context, messages []analytics.Event) error {
	queryBuilder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Insert("analytics_events").
		Columns(
			"tg_user_id",
			"event_name",
			"event_timestamp",
		)

	for _, m := range messages {
		queryBuilder = queryBuilder.Values(
			m.TgUserID,
			m.EventName,
			m.Timestamp,
		)
	}

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return errors.Wrap(err, "creating insert query")
	}

	_, err = c.db.ExecContext(ctx, query, args...)
	if err != nil {
		return errors.Wrap(err, "executing insert query")
	}

	return nil
}

func (c *CH) Write(message analytics.Event) {
	go func() {
		c.messageQueue <- message
	}()
}

func (c *CH) StartWorker(ctx context.Context) {
	go func() {
		batch := makeNewBlancBatch()
		for {
			select {
			case m := <-c.messageQueue:
				batch = append(batch, m)
			case <-time.After(1 * time.Second):
				c.handleBatch(ctx, batch)
				batch = makeNewBlancBatch()
				continue
			case <-ctx.Done():
				c.handleBatch(ctx, batch)
				return
			}
		}
	}()
}

func makeNewBlancBatch() []analytics.Event {
	return make([]analytics.Event, 0, 100)
}

func (c *CH) handleBatch(ctx context.Context, batch []analytics.Event) {
	if len(batch) != 0 {
		err := c.WriteToDB(ctx, batch)
		if err != nil {
			c.logger.Error(err.Error())
		}
	}
}
