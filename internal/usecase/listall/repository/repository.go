package repository

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/gwassel/TasksOfWoe/internal/domain"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type repository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *repository {
	return &repository{db: db}
}

func (r *repository) ListAllTasks(userID int64) ([]domain.Task, error) {
	op := "list all tasks"

	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select(
			"*",
		).
		From("tasks").
		Where(sq.And{sq.Eq{"user_id": userID}}).
		OrderBy("completed_at NULLS FIRST", "is_in_work DESC", "user_task_id ASC")

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to build query '%s'", query)
	}
	query = fmt.Sprintf("-- %s\n%s", op, query)

	var tasks []domain.Task
	err = r.db.SelectContext(context.TODO(), &tasks, query, args...)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to exec query '%s'", query)
	}

	return tasks, nil
}
