package domain

type URL struct {
	ID          string
	OriginalURL string
	UserID      string
	IsDeleted   bool
}

func NewURL(id, originalURL string, userID string, isDeleted bool) *URL {
	return &URL{
		ID:          id,
		OriginalURL: originalURL,
		UserID:      userID,
		IsDeleted:   isDeleted,
	}
}
