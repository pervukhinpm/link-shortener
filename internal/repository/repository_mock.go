package repository

import (
	"fmt"
	"github.com/pervukhinpm/link-shortener.git/domain"
)

type MockRepository struct {
	Urls map[string]*domain.URL
}

func NewMockRepository() *MockRepository {
	return &MockRepository{
		Urls: make(map[string]*domain.URL),
	}
}

func (m *MockRepository) Add(url *domain.URL) error {
	m.Urls[url.ID] = url
	return nil
}

func (m *MockRepository) Get(id string) (*domain.URL, error) {
	url, exists := m.Urls[id]
	if !exists {
		return nil, fmt.Errorf("URL not found")
	}
	return url, nil
}
