package repository

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

// NewRepository создает новый экземпляр репозитория в зависимости от конфигурации.
// Приоритет выбора:
// 1. DatabaseRepository, если указан DSN и есть подключение к БД
// 2. FileRepository, если указан путь к файловому хранилищу
// 3. RAMRepository в остальных случаях
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
