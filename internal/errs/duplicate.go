package errs

import (
	"fmt"
	"github.com/pervukhinpm/link-shortener.git/domain"
)

type OriginalURLAlreadyExists struct {
	URL *domain.URL
}

func NewOriginalURLAlreadyExists(url *domain.URL) *OriginalURLAlreadyExists {
	return &OriginalURLAlreadyExists{URL: url}
}

func (e *OriginalURLAlreadyExists) Error() string {
	return fmt.Sprintf("original URL already exists: %s", e.URL.OriginalURL)
}
