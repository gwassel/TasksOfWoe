package space

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

// TODO get or create
func (r *repository) GetUserByTgId(ctx context.Context, chatID int64) (*int64, error) {
	op := "get space by tg id"

	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select(
			"id",
		).
		From("user").
		Where(sq.And{sq.Eq{"external_source": user_domain.ExternalSourceTG}, sq.Eq{"external_user_id": chatID}})

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to build query '%s'", query)
	}
	query = fmt.Sprintf("-- %s\n%s", op, query)

	var result int64
	err = r.db.SelectContext(ctx, &result, query, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.Wrapf(err, "failed to exec query '%s'", query)
	}

	return &result, nil
}
