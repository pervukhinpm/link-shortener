// Package service предоставляет бизнес-логику приложения.
// Включает в себя:
//   - Создание коротких URL
//   - Управление URL пользователей
//   - Проверку состояния сервиса
package service

import (
	"context"
	"errors"

	"github.com/pervukhinpm/link-shortener.git/domain"
	"github.com/pervukhinpm/link-shortener.git/internal/model"
)

// MockShortenerService представляет мок-реализацию сервиса для работы с URL.
// Используется для тестирования без реальной базы данных.
type MockShortenerService struct {
	ShortenURL *domain.URL
}

// NewMockService создает новый экземпляр MockShortenerService.
// Инициализирует пустую структуру для тестирования.
func NewMockService() *MockShortenerService {
	return &MockShortenerService{}
}

// Shorten создает короткий URL из оригинального.
// Возвращает предустановленный URL или ошибку, если сервис не инициализирован.
func (u *MockShortenerService) Shorten(original string, ctx context.Context) (*domain.URL, error) {
	if u.ShortenURL == nil {
		return nil, errors.New("shorten service not found")
	}
	return u.ShortenURL, nil
}

// Find ищет URL по его короткому идентификатору.
// Возвращает предустановленный URL или ошибку, если сервис не инициализирован.
func (u *MockShortenerService) Find(id string, ctx context.Context) (*domain.URL, error) {
	if u.ShortenURL == nil {
		return nil, errors.New("shorten service not found")
	}
	return u.ShortenURL, nil
}

// AddBatch добавляет несколько URL в хранилище.
// Обрабатывает каждый URL через метод Shorten.
func (u *MockShortenerService) AddBatch(urls []domain.URL, ctx context.Context) error {
	for _, url := range urls {
		if _, err := u.Shorten(url.OriginalURL, ctx); err != nil {
			return err
		}
	}
	return nil
}

// GetByUserID возвращает все URL, созданные пользователем.
// Если сервис не инициализирован, возвращает пустой слайс.
func (u *MockShortenerService) GetByUserID(ctx context.Context) (*[]domain.URL, error) {
	if u.ShortenURL == nil {
		emptyURLs := make([]domain.URL, 0)
		return &emptyURLs, nil
	}
	urls := []domain.URL{*u.ShortenURL}
	return &urls, nil
}

// GetFlagByShortURL проверяет, был ли URL удален.
// В мок-реализации всегда возвращает false.
func (u *MockShortenerService) GetFlagByShortURL(ctx context.Context, shortURL string) (bool, error) {
	return false, nil
}

// DeleteURLBatch помечает несколько URL как удаленные.
// В мок-реализации не выполняет никаких действий.
func (u *MockShortenerService) DeleteURLBatch(ctx context.Context, deleteBatch model.DeleteBatch) {

}
