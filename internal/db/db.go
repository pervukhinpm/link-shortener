package db

import (
	"database/sql"
	"github.com/pervukhinpm/link-shortener.git/internal/middleware"
)

func NewDB(dsn string) (*sql.DB, error) {
	pool, err := sql.Open("pgx", dsn)
	if err != nil {
		middleware.Log.Error("Failed to connect to database: %v", err)
		return nil, err
	}
	return pool, nil
}
