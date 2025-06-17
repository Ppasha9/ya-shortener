package model

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponse struct {
	Result string `json:"result"`
}

type ShortenBatchItemRequest struct {
	CorrID  string `json:"correlation_id"`
	OrigURL string `json:"original_url"`
}

type ShortenBatchItemResponse struct {
	CorrID   string `json:"correlation_id"`
	ShortURL string `json:"short_url"`
}

type URLsPair struct {
	ShortURL string
	OrigURL  string
}
