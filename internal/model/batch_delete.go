package model

// DeleteBatch представляет запрос на пакетное удаление URL.
// Содержит информацию о пользователе и списке URL для удаления.
type DeleteBatch struct {
	// UserID - идентификатор пользователя
	UserID string `json:"user_id"`
	// ShortenedURL - список коротких идентификаторов URL для удаления
	ShortenedURL []string `json:"shortened_url"`
}
