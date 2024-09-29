package repository

import (
	"context"
	"errors"
	"github.com/pervukhinpm/link-shortener.git/domain"
	"github.com/pervukhinpm/link-shortener.git/internal/errs"
)

type RAMRepository struct {
	MapURL map[string]string
}

func NewRAMRepository() (*RAMRepository, error) {
	return &RAMRepository{MapURL: make(map[string]string)}, nil
}

func (rmr *RAMRepository) Add(url *domain.URL, ctx context.Context) error {
	for _, existingOriginalURL := range rmr.MapURL {
		if existingOriginalURL == url.OriginalURL {
			return errs.NewOriginalURLAlreadyExists(url)
		}
	}

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

func (rmr *RAMRepository) AddBatch(urls []domain.URL, ctx context.Context) error {
	for _, url := range urls {
		if err := rmr.Add(&url, ctx); err != nil {
			return err
		}
	}
	return nil
}

func (rmr *RAMRepository) Close() error {
	return nil
}
