package repository

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type repository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *repository {
	return &repository{db: db}
}

func (r *repository) AddTask(userID int64, task string) (int64, error) {
	op := "insert task"

	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Insert("tasks").
		Columns("user_id", "task", "user_task_id").
		Values(userID, task, sq.Expr("(SELECT COALESCE(MAX(user_task_id)+1, 1) FROM tasks where user_id=$3)", userID)).
		Suffix("RETURNING user_task_id")

	query, args, err := builder.ToSql()
	if err != nil {
		return 0, errors.Wrapf(err, "failed to build query '%s'", query)
	}
	query = fmt.Sprintf("-- %s\n%s", op, query)

	result, err := r.db.QueryContext(context.TODO(), query, args...)
	if err != nil {
		return 0, errors.Wrapf(err, "failed to exec query '%s'", query)
	}
	defer func() { _ = result.Close() }()

	var userTaskID int64
	for result.Next() {
		err = result.Scan(&userTaskID)
		if err != nil {
			return 0, errors.Wrapf(err, "failed to retrieve task ID")
		}
	}

	return userTaskID, nil
}
