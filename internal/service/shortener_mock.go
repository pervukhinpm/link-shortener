package service

import (
	"context"
	"errors"
	"github.com/pervukhinpm/link-shortener.git/domain"
)

type MockShortenerService struct {
	ShortenURL *domain.URL
}

func NewMockService() *MockShortenerService {
	return &MockShortenerService{}
}

func (u *MockShortenerService) Shorten(original string, ctx context.Context) (*domain.URL, error) {
	if u.ShortenURL == nil {
		return nil, errors.New("shorten service not found")
	}
	return u.ShortenURL, nil
}

func (u *MockShortenerService) Find(id string, ctx context.Context) (*domain.URL, error) {
	if u.ShortenURL == nil {
		return nil, errors.New("shorten service not found")
	}
	return u.ShortenURL, nil
}

func (u *MockShortenerService) AddBatch(urls []domain.URL, ctx context.Context) error {
	for _, url := range urls {
		if _, err := u.Shorten(url.OriginalURL, ctx); err != nil {
			return err
		}
	}
	return nil
}

func (u *MockShortenerService) GetByUserID(ctx context.Context) (*[]domain.URL, error) {
	if u.ShortenURL == nil {
		return nil, errors.New("shorten service not found")
	}
	urls := []domain.URL{*u.ShortenURL}
	return &urls, nil
}
