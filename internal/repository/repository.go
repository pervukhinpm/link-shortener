package repository

import (
	"context"
	"github.com/pervukhinpm/link-shortener.git/domain"
)

type Repository interface {
	Add(url *domain.URL, ctx context.Context) error
	AddBatch(urls []domain.URL, ctx context.Context) error
	Get(id string, ctx context.Context) (*domain.URL, error)
	GetByUserID(ctx context.Context) (*[]domain.URL, error)
	GetFlagByShortURL(ctx context.Context, shortenedURL string) (bool, error)
	DeleteURLBatch(ctx context.Context, urls []UserShortURL) error
	Close() error
}
