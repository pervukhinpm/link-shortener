package repository

import (
	"context"
	"errors"

	"github.com/pervukhinpm/link-shortener.git/domain"
	"github.com/pervukhinpm/link-shortener.git/internal/errs"
	"github.com/pervukhinpm/link-shortener.git/internal/middleware"
)

// RAMRepository реализует интерфейс Repository для работы с оперативной памятью.
// Использует map для хранения данных.
type RAMRepository struct {
	MapURL map[string]domain.URL
}

// NewRAMRepository создает новый экземпляр RAMRepository.
// Инициализирует пустое хранилище в памяти.
func NewRAMRepository() (*RAMRepository, error) {
	return &RAMRepository{MapURL: make(map[string]domain.URL)}, nil
}

// Add добавляет новый URL в хранилище в памяти.
// Проверяет уникальность оригинального URL.
func (rmr *RAMRepository) Add(url *domain.URL, ctx context.Context) error {
	for _, existingURL := range rmr.MapURL {
		if existingURL.OriginalURL == url.OriginalURL {
			return errs.NewOriginalURLAlreadyExists(url)
		}
	}

	rmr.MapURL[url.ID] = *url
	return nil
}

// Get возвращает URL по его короткому идентификатору.
// Если URL не найден, возвращает ошибку.
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

// AddBatch добавляет несколько URL в хранилище в памяти.
// Вызывает Add для каждого URL в пакете.
func (rmr *RAMRepository) AddBatch(urls []domain.URL, ctx context.Context) error {
	for _, url := range urls {
		if err := rmr.Add(&url, ctx); err != nil {
			return err
		}
	}
	return nil
}

// Close закрывает хранилище в памяти.
// В данном случае просто возвращает nil, так как нет ресурсов для освобождения.
func (rmr *RAMRepository) Close() error {
	return nil
}

// GetByUserID возвращает все URL, созданные пользователем.
// Если URL не найдены, возвращает пустой слайс.
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

	// Возвращаем список URL (может быть пустым)
	return &urls, nil
}

// GetFlagByShortURL проверяет, был ли URL удален.
// Если URL не найден, возвращает ошибку ErrURLNotFound.
func (rmr *RAMRepository) GetFlagByShortURL(ctx context.Context, shortenedURL string) (bool, error) {
	urlData, exists := rmr.MapURL[shortenedURL]
	if !exists {
		return false, errs.ErrURLNotFound
	}

	// Возвращаем флаг is_deleted
	return urlData.IsDeleted, nil
}

// DeleteURLBatch помечает несколько URL как удаленные.
// Проверяет принадлежность URL пользователю перед удалением.
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
