package repository

import (
	"context"
	"errors"
	"github.com/pervukhinpm/link-shortener.git/domain"
)

type RAMRepository struct {
	MapURL map[string]string
}

func NewRAMRepository() (*RAMRepository, error) {
	return &RAMRepository{MapURL: make(map[string]string)}, nil
}

func (rmr *RAMRepository) Add(url *domain.URL, ctx context.Context) error {
	rmr.MapURL[url.ID] = url.OriginalURL
	return nil
}

func (rmr *RAMRepository) Get(id string, ctx context.Context) (*domain.URL, error) {
	longURL := rmr.MapURL[id]
	if longURL == "" {
		return nil, errors.New("url not found")
	}
	return domain.NewURL(id, longURL), nil
}

func (rmr *RAMRepository) Close() error {
	return nil
}
