package repository

import (
	"github.com/pervukhinpm/link-shortener.git/domain"
)

type Repository interface {
	Add(url *domain.URL) error
	Get(id string) (*domain.URL, error)
}
