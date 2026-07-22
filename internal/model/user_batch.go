package model

// URLByUserBatchResponseItem представляет элемент ответа на запрос получения URL пользователя.
// Содержит информацию о сокращенном и оригинальном URL.
type URLByUserBatchResponseItem struct {
	// ShortURL - сокращенный URL
	ShortURL string `json:"short_url"`
	// OriginalURL - оригинальный URL
	OriginalURL string `json:"original_url"`
}
