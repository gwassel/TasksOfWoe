package performance

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type AdminRepository interface {
	GetAdminUsers(ctx context.Context) ([]int64, error)
	AddAdminUser(ctx context.Context, telegramUserID int64) error
	RemoveAdminUser(ctx context.Context, telegramUserID int64) error
	IsAdminUser(ctx context.Context, telegramUserID int64) (bool, error)
}

type adminRepository struct {
	db *sqlx.DB
}

func NewAdminRepository(db *sqlx.DB) AdminRepository {
	return &adminRepository{db: db}
}

func (r *adminRepository) GetAdminUsers(ctx context.Context) ([]int64, error) {
	query := `SELECT telegram_user_id FROM admin_users ORDER BY registered_at`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, errors.Wrap(err, "getting admin users")
	}
	defer rows.Close()

	var adminIDs []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, errors.Wrap(err, "scanning admin user id")
		}
		adminIDs = append(adminIDs, id)
	}

	return adminIDs, nil
}

func (r *adminRepository) AddAdminUser(ctx context.Context, telegramUserID int64) error {
	query := `INSERT INTO admin_users (telegram_user_id) VALUES ($1) ON CONFLICT (telegram_user_id) DO NOTHING`

	_, err := r.db.ExecContext(ctx, query, telegramUserID)
	if err != nil {
		return errors.Wrap(err, "adding admin user")
	}

	return nil
}

func (r *adminRepository) RemoveAdminUser(ctx context.Context, telegramUserID int64) error {
	query := `DELETE FROM admin_users WHERE telegram_user_id = $1`

	result, err := r.db.ExecContext(ctx, query, telegramUserID)
	if err != nil {
		return errors.Wrap(err, "removing admin user")
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("admin user not found")
	}

	return nil
}

func (r *adminRepository) IsAdminUser(ctx context.Context, telegramUserID int64) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM admin_users WHERE telegram_user_id = $1)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, telegramUserID).Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "checking admin user")
	}

	return exists, nil
}
