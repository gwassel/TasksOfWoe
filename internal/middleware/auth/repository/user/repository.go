package user

import (
	"context"
	"database/sql"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	user_domain "github.com/gwassel/TasksOfWoe/internal/domain/user"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type repository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *repository {
	return &repository{
		db: db,
	}
}

func (r *repository) GetOrCreate(ctx context.Context, userTgID int64) (*int64, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to begin transaction")
	}

	selectBuilder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select(
			"id",
		).
		From("user").
		Where(sq.And{sq.Eq{"external_source": user_domain.ExternalSourceTG}, sq.Eq{"external_user_id": userTgID}})

	query, args, err := selectBuilder.ToSql()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to build query '%s'", query)
	}
	query = fmt.Sprintf("-- %s\n%s", "get user by tg id", query)

	var result int64
	err = tx.GetContext(ctx, &result, query, args...)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, errors.Wrapf(err, "failed to execute select query '%s'", query)
		}

		if err := tx.Commit(); err != nil {
			return nil, errors.Wrap(err, "failed to commit transaction")
		}
		return &result, nil
	}

	insertBuilder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Insert("user").
		Columns("external_source", "external_user_id").
		Values(user_domain.ExternalSourceTG, userTgID).
		Suffix("RETURNING id")

	query, args, err = insertBuilder.ToSql()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to build insert query '%s'", query)
	}
	query = fmt.Sprintf("-- %s\n%s", "create user", query)

	if err := tx.GetContext(ctx, &result, query, args...); err != nil {
		return nil, errors.Wrapf(err, "failed to execute insert query '%s'", query)
	}

	if err := tx.Commit(); err != nil {
		return nil, errors.Wrap(err, "failed to commit transaction")
	}

	return &result, nil
}
