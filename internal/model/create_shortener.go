package model

type CreateShortenerBody struct {
	URL string `json:"service"`
}

type CreateShortenerResponse struct {
	Result string `json:"result"`
}
