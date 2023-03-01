package responses

type ResponseShortLink struct {
	Result string `json:"result"`
}

type ResponseShortOriginalLink struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type ResponseLinksBatch struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}
