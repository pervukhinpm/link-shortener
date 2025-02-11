package repository

import (
	"context"
	"errors"
	"github.com/pervukhinpm/link-shortener.git/domain"
	"github.com/pervukhinpm/link-shortener.git/internal/errs"
	"github.com/pervukhinpm/link-shortener.git/internal/middleware"
)

type RAMRepository struct {
	MapURL map[string]domain.URL
}

func NewRAMRepository() (*RAMRepository, error) {
	return &RAMRepository{MapURL: make(map[string]domain.URL)}, nil
}

func (rmr *RAMRepository) Add(url *domain.URL, ctx context.Context) error {
	for _, existingURL := range rmr.MapURL {
		if existingURL.OriginalURL == url.OriginalURL {
			return errs.NewOriginalURLAlreadyExists(url)
		}
	}

	rmr.MapURL[url.ID] = *url
	return nil
}

func (rmr *RAMRepository) Get(id string, ctx context.Context) (*domain.URL, error) {
	longURL := rmr.MapURL[id].OriginalURL
	userID := middleware.GetUserID(ctx)
	isDeleted := rmr.MapURL[id].IsDeleted
	if longURL == "" {
		return nil, errors.New("url not found")
	}
	url := domain.NewURL(id, longURL, userID, isDeleted)
	return url, nil
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

func (rmr *RAMRepository) GetByUserID(ctx context.Context) (*[]domain.URL, error) {
	var urls []domain.URL

	// Получаем текущий UserID из контекста
	userID := middleware.GetUserID(ctx)

	// Проходим по всем URL в хранилище
	for _, url := range rmr.MapURL {
		// Сравниваем UserID
		if url.UserID == userID {
			urls = append(urls, url)
		}
	}

	// Если не найдено ни одной записи
	if len(urls) == 0 {
		return nil, errors.New("no urls found for this user")
	}

	// Возвращаем список URL
	return &urls, nil
}

func (rmr *RAMRepository) GetFlagByShortURL(ctx context.Context, shortenedURL string) (bool, error) {
	urlData, exists := rmr.MapURL[shortenedURL]
	if !exists {
		return false, errs.ErrURLNotFound
	}

	// Возвращаем флаг is_deleted
	return urlData.IsDeleted, nil
}

func (rmr *RAMRepository) DeleteURLBatch(ctx context.Context, urls []UserShortURL) error {
	for _, url := range urls {
		urlData, exists := rmr.MapURL[url.ShortURL]
		if !exists {
			continue
		}
		if urlData.UserID == url.UserID {
			urlData.IsDeleted = true
			rmr.MapURL[url.ShortURL] = urlData
		}
	}

	return nil
}
