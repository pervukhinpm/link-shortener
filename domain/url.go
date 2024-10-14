package domain

type URL struct {
	ID          string
	OriginalURL string
	UserID      string
	IsDeleted   bool
}

func NewURL(id, originalURL string, userID string, IsDeleted bool) *URL {
	return &URL{id, originalURL, userID, IsDeleted}
}
