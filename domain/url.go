package domain

// URL представляет собой сокращенный URL в системе.
// Содержит информацию об оригинальном URL, коротком идентификаторе и статусе.
type URL struct {
	// ID - короткий идентификатор URL
	ID string `json:"id"`
	// OriginalURL - оригинальный URL
	OriginalURL string `json:"original_url"`
	// UserID - идентификатор пользователя, создавшего URL
	UserID string `json:"user_id"`
	// IsDeleted - флаг, указывающий, был ли URL удален
	IsDeleted bool `json:"is_deleted"`
}

// NewURL создает новый экземпляр URL с указанными параметрами.
func NewURL(id, originalURL, userID string, isDeleted bool) *URL {
	return &URL{
		ID:          id,
		OriginalURL: originalURL,
		UserID:      userID,
		IsDeleted:   isDeleted,
	}
}
