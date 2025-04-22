package model

// BatchRequestItem представляет элемент запроса на пакетное создание URL.
// Содержит информацию об оригинальном URL и его корреляционном идентификаторе.
type BatchRequestItem struct {
	// CorrelationID - идентификатор для сопоставления запроса и ответа
	CorrelationID string `json:"correlation_id"`
	// OriginalURL - оригинальный URL для сокращения
	OriginalURL string `json:"original_url"`
}

// BatchRequestBody представляет тело запроса на пакетное создание URL.
// Содержит массив URL для сокращения.
type BatchRequestBody struct {
	// BatchList - список URL для сокращения
	BatchList []BatchRequestItem `json:"batch"`
}

// BatchResponseItem представляет элемент ответа на пакетное создание URL.
// Содержит информацию о созданном сокращенном URL и его корреляционном идентификаторе.
type BatchResponseItem struct {
	// CorrelationID - идентификатор для сопоставления запроса и ответа
	CorrelationID string `json:"correlation_id"`
	// ShortURL - сокращенный URL
	ShortURL string `json:"short_url"`
}
