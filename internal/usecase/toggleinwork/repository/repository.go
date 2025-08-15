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

func (r *repository) ToggleInWork(userID int64, userTaskIDs []int64, value bool) error {
	op := "update inWork"

	var taken_at sq.Sqlizer
	if value {
		taken_at = sq.Expr("NOW()")
	} else {
		taken_at = sq.Expr("NULL")
	}
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Update("tasks").
		Set("is_in_work", value).
		Set("taken_at", taken_at).
		Where(sq.And{sq.Eq{"user_id": userID}, sq.Eq{"user_task_id": userTaskIDs}, sq.Eq{"completed": false}, sq.Eq{"is_in_work": !value}})

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
