package repository

import (
	"context"
	"fmt"
	"sync"

	"github.com/pervukhinpm/link-shortener.git/domain"
)

type MockRepository struct {
	Urls map[string]*domain.URL
	mu   sync.RWMutex
}

func NewMockRepository() *MockRepository {
	return &MockRepository{
		Urls: make(map[string]*domain.URL),
	}
}

func (m *MockRepository) Add(url *domain.URL, ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Urls[url.ID] = url
	return nil
}

func (m *MockRepository) Get(id string, ctx context.Context) (*domain.URL, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	url, exists := m.Urls[id]
	if !exists {
		return nil, fmt.Errorf("URL not found")
	}
	return url, nil
}

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

func (m *MockRepository) GetByUserID(ctx context.Context) (*[]domain.URL, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make([]domain.URL, 0, len(m.Urls))
	for _, url := range m.Urls {
		result = append(result, *url)
	}
	return &result, nil
}

func (m *MockRepository) GetFlagByShortURL(ctx context.Context, shortenedURL string) (bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	url, exists := m.Urls[shortenedURL]
	if !exists {
		return false, fmt.Errorf("URL not found")
	}
	return url.IsDeleted, nil
}

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

func (m *MockRepository) Close() error {
	return nil
}
