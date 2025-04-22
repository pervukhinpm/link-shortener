// Package errs предоставляет пользовательские ошибки для приложения.
// Включает в себя:
//   - Ошибки для работы с URL
//   - Ошибки для работы с базой данных
package errs

import (
	"fmt"

	"github.com/pervukhinpm/link-shortener.git/domain"
)

// OriginalURLAlreadyExists представляет ошибку, возникающую при попытке добавить уже существующий URL.
// Содержит информацию о существующем URL.
type OriginalURLAlreadyExists struct {
	// URL - существующий URL
	URL *domain.URL
}

// NewOriginalURLAlreadyExists создает новый экземпляр ошибки OriginalURLAlreadyExists.
// Принимает URL, который уже существует в базе данных.
func NewOriginalURLAlreadyExists(url *domain.URL) *OriginalURLAlreadyExists {
	return &OriginalURLAlreadyExists{URL: url}
}

// Error возвращает строковое представление ошибки.
// Реализует интерфейс error.
func (e *OriginalURLAlreadyExists) Error() string {
	return fmt.Sprintf("URL %s already exists", e.URL.OriginalURL)
}
