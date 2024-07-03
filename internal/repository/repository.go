package repository

import (
	"errors"
	"github.com/pervukhinpm/link-shortener.git/domain"
)

type Repository interface {
	Add(url *domain.URL) error
	Get(id string) (*domain.URL, error)
}

type InMemoryRepository struct {
	storage map[string]string
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{storage: make(map[string]string)}
}

func (r *InMemoryRepository) Add(url *domain.URL) error {
	r.storage[url.ID] = url.OriginalURL
	return nil
}

func (r *InMemoryRepository) Get(id string) (*domain.URL, error) {
	url, exists := r.storage[id]
	if !exists {
		return nil, errors.New("URL not found")
	}
	return domain.NewURL(id, url), nil
}
