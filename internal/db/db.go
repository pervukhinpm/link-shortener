package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pervukhinpm/link-shortener.git/internal/middleware"
)

// NewDB создает новый пул соединений с базой данных PostgreSQL.
// Использует указанный DSN для подключения.
func NewDB(dsn string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		middleware.Log.Error("Failed to connect to database: %v", err)
		return nil, err
	}
	return pool, nil
}
