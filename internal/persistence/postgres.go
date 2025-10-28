package persistence

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
)

func NewDB() (*sqlx.DB, error) {
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	connStr := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s port=%s",
		dbHost, dbUser, dbName, dbPassword, dbPort)
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		return nil, err
	}
	return db, nil
}
