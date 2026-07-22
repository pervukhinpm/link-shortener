package model

// CreateShortenerBody представляет тело запроса на создание сокращенного URL.
// Содержит оригинальный URL для сокращения.
type CreateShortenerBody struct {
	// URL - оригинальный URL для сокращения
	URL string `json:"url"`
}

// CreateShortenerResponse представляет ответ на запрос создания сокращенного URL.
// Содержит созданный сокращенный URL.
type CreateShortenerResponse struct {
	// Result - созданный сокращенный URL
	Result string `json:"result"`
}
