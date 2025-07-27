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

func (r *repository) CompleteTask(userID int64, userTaskID int64) error {
	op := "complete task"

	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Update("tasks").
		Set("completed", true).
		Set("completed_at", sq.Expr("NOW()")).
		Where(sq.And{sq.Eq{"user_task_id": userTaskID}, sq.Eq{"user_id": userID}, sq.Eq{"completed": false}})

	query, args, err := builder.ToSql()
	if err != nil {
		return errors.Wrapf(err, "failed to build query '%s'", query)
	}
	query = fmt.Sprintf("-- %s\n%s", op, query)

	_, err = r.db.ExecContext(context.TODO(), query, args...)
	if err != nil {
		return errors.Wrapf(err, "failed to exec query '%s'", query)
	}
	return err

}
