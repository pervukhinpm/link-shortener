package url

import (
	"errors"
	"github.com/pervukhinpm/link-shortener.git/domain"
)

type MockShortenerService struct {
	ShortenUrl *domain.URL
}

func NewMockService() *MockShortenerService {
	return &MockShortenerService{}
}

func (u *MockShortenerService) Shorten(original string) (*domain.URL, error) {
	if u.ShortenUrl == nil {
		return nil, errors.New("shorten url not found")
	}
	return u.ShortenUrl, nil
}

func (u *MockShortenerService) Find(id string) (*domain.URL, error) {
	if u.ShortenUrl == nil {
		return nil, errors.New("shorten url not found")
	}
	return u.ShortenUrl, nil
}
