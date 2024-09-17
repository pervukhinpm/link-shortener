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
