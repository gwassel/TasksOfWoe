package repository

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	domain "github.com/gwassel/TasksOfWoe/internal/domain/task"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type repository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *repository {
	return &repository{db: db}
}

func (r *repository) GetUnencryptedTasks() ([]domain.Task, error) {
	op := "Get unencrypted tasks"

	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select("*").
		From("tasks").
		Where(sq.Eq{"task": nil})

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

func (r *repository) EncryptTask(taskID int64, task []byte) error {
	op := "Encode task"

	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Update("tasks").
		Set("encrypted_task", task).
		Where(sq.Eq{"id": taskID})

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
