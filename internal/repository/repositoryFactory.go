package repository

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewRepository(
	dsn string,
	fileStoragePath string,
	db *pgxpool.Pool,
) (Repository, error) {
	// Если есть DSN и подключение к БД, создаем DatabaseRepository
	if dsn != "" && db != nil {
		return NewDatabaseRepository(db)
	}

	// Если есть путь к файловому хранилищу, создаем FileRepository
	if fileStoragePath != "" {
		return NewFileRepository(fileStoragePath)
	}

	return NewRAMRepository()
}
