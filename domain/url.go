package domain

type URL struct {
	ID          string
	OriginalURL string
}

func NewURL(id, originalURL string) *URL {
	return &URL{id, originalURL}
}
