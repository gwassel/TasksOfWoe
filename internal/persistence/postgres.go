package persistence

import (
	"github.com/jmoiron/sqlx"
)

func NewDB(connStr string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		return nil, err
	}
	return db, nil
}
