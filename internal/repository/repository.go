package repository

import (
	"context"
	"github.com/pervukhinpm/link-shortener.git/domain"
)

type Repository interface {
	Add(url *domain.URL, ctx context.Context) error
	AddBatch(urls []domain.URL, ctx context.Context) error
	Get(id string, ctx context.Context) (*domain.URL, error)
	Close() error
}
