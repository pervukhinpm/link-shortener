package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"github.com/pervukhinpm/link-shortener.git/domain"
	"github.com/pervukhinpm/link-shortener.git/internal/repository"
	"strings"
)

type ShortenerServiceReaderWriter interface {
	Find(id string, ctx context.Context) (*domain.URL, error)
	AddBatch(urls []domain.URL, ctx context.Context) error
	Shorten(original string, ctx context.Context) (*domain.URL, error)
}

type ShortenerService struct {
	repo repository.Repository
}

func NewURLService(repo repository.Repository) *ShortenerService {
	return &ShortenerService{repo: repo}
}

func (u *ShortenerService) Shorten(original string, ctx context.Context) (*domain.URL, error) {
	randomBytes := make([]byte, 6)
	if _, err := rand.Read(randomBytes); err != nil {
		return nil, err
	}
	short := base64.URLEncoding.EncodeToString(randomBytes)
	short = strings.TrimRight(short, "=")
	url := domain.NewURL(short, original)
	if err := u.repo.Add(url, ctx); err != nil {
		return nil, err
	}
	return url, nil
}

func (u *ShortenerService) AddBatch(urls []domain.URL, ctx context.Context) error {
	if err := u.repo.AddBatch(urls, ctx); err != nil {
		return err
	}
	return nil
}

func (u *ShortenerService) Find(id string, ctx context.Context) (*domain.URL, error) {
	url, err := u.repo.Get(id, ctx)
	if err != nil {
		return nil, err
	}
	return url, nil
}
