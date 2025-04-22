// Package repository предоставляет реализацию интерфейса Repository для работы с базой данных.
// Включает в себя:
//   - Работу с PostgreSQL
//   - Управление соединениями
//   - Операции с URL
package repository

import (
	"context"
	"fmt"
	"sync"

	"github.com/pervukhinpm/link-shortener.git/domain"
)

// MockRepository представляет мок-реализацию интерфейса Repository.
// Используется для тестирования без реальной базы данных.
type MockRepository struct {
	// Urls - хранилище URL в памяти
	Urls map[string]*domain.URL
	mu   sync.RWMutex
}

// NewMockRepository создает новый экземпляр MockRepository.
// Инициализирует хранилище URL в памяти.
func NewMockRepository() *MockRepository {
	return &MockRepository{
		Urls: make(map[string]*domain.URL),
	}
}

// Add добавляет новый URL в хранилище.
// Генерирует UUID для записи и проверяет уникальность оригинального URL.
func (m *MockRepository) Add(url *domain.URL, ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Urls[url.ID] = url
	return nil
}

// Get возвращает URL по его короткому идентификатору.
// Если URL не найден, возвращает ошибку.
func (m *MockRepository) Get(id string, ctx context.Context) (*domain.URL, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	url, exists := m.Urls[id]
	if !exists {
		return nil, fmt.Errorf("URL not found")
	}
	return url, nil
}

// AddBatch добавляет несколько URL в хранилище в рамках одной транзакции.
// Использует пакетную вставку для оптимизации производительности.
func (m *MockRepository) AddBatch(urls []domain.URL, ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if len(m.Urls) == 0 {
		m.Urls = make(map[string]*domain.URL, len(urls))
	}

	for i := range urls {
		url := urls[i]
		m.Urls[url.ID] = &url
	}
	return nil
}

// GetByUserID возвращает все URL, созданные пользователем.
// Если URL не найдены, возвращает пустой слайс.
func (m *MockRepository) GetByUserID(ctx context.Context) (*[]domain.URL, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make([]domain.URL, 0, len(m.Urls))
	for _, url := range m.Urls {
		result = append(result, *url)
	}
	return &result, nil
}

// GetFlagByShortURL проверяет, был ли URL удален.
// Если URL не найден, возвращает ошибку ErrURLNotFound.
func (m *MockRepository) GetFlagByShortURL(ctx context.Context, shortenedURL string) (bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	url, exists := m.Urls[shortenedURL]
	if !exists {
		return false, fmt.Errorf("URL not found")
	}
	return url.IsDeleted, nil
}

// DeleteURLBatch помечает несколько URL как удаленные.
// Использует пакетное обновление для оптимизации производительности.
func (m *MockRepository) DeleteURLBatch(ctx context.Context, urls []UserShortURL) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, url := range urls {
		if storedURL, exists := m.Urls[url.ShortURL]; exists {
			storedURL.IsDeleted = true
		}
	}
	return nil
}

// Close закрывает соединение с хранилищем.
func (m *MockRepository) Close() error {
	return nil
}
