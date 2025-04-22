package repository

import (
	"context"

	"github.com/pervukhinpm/link-shortener.git/domain"
)

// Repository определяет интерфейс для работы с хранилищем URL.
// Предоставляет методы для добавления, получения и удаления URL.
type Repository interface {
	// Add добавляет новый URL в хранилище.
	Add(url *domain.URL, ctx context.Context) error
	// AddBatch добавляет несколько URL в хранилище.
	AddBatch(urls []domain.URL, ctx context.Context) error
	// Get возвращает URL по его короткому идентификатору.
	Get(id string, ctx context.Context) (*domain.URL, error)
	// GetByUserID возвращает все URL, созданные пользователем.
	GetByUserID(ctx context.Context) (*[]domain.URL, error)
	// GetFlagByShortURL проверяет, был ли URL удален.
	GetFlagByShortURL(ctx context.Context, shortenedURL string) (bool, error)
	// DeleteURLBatch удаляет несколько URL пользователя.
	DeleteURLBatch(ctx context.Context, urls []UserShortURL) error
	// Close закрывает соединение с хранилищем.
	Close() error
}
