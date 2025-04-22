package repository

// UserShortURL представляет связь между пользователем и его сокращенным URL.
// Используется для пакетного удаления URL.
type UserShortURL struct {
	// UserID идентификатор пользователя
	UserID string
	// ShortURL короткий идентификатор URL
	ShortURL string
}
