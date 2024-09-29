package model

type BatchRequestBody struct {
	BatchList []BatchRequestBodyItem
}

type BatchRequestBodyItem struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type BatchResponse struct {
	BatchList []BatchResponseItem
}

type BatchResponseItem struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}
