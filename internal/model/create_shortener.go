package model

type CreateShortenerBody struct {
	URL string `json:"url"`
}

type CreateShortenerResponse struct {
	Result string `json:"result"`
}
