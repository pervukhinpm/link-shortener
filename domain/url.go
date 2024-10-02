package domain

type URL struct {
	ID          string
	OriginalURL string
	UserID      string
}

func NewURL(id, originalURL string, userID string) *URL {
	return &URL{id, originalURL, userID}
}
