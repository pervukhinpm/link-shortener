package url

import (
	"crypto/sha1"
	"encoding/hex"
	"github.com/pervukhinpm/link-shortener.git/domain"
	"github.com/pervukhinpm/link-shortener.git/internal/repository"
)

type Service struct {
	repo repository.Repository
}

func NewURLService(repo repository.Repository) *Service {
	return &Service{repo: repo}
}

func (u *Service) Shorten(original string) (*domain.URL, error) {
	hash := sha1.New()
	hash.Write([]byte(original))
	short := hex.EncodeToString(hash.Sum(nil))[:8]
	url := domain.NewURL(short, original)
	if err := u.repo.Add(url); err != nil {
		return nil, err
	}
	return url, nil
}

func (u *Service) Find(id string) (*domain.URL, error) {
	url, err := u.repo.Get(id)
	if err != nil {
		return nil, err
	}
	return url, nil
}
