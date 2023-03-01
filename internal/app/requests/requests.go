package requests

type RequestLink struct {
	URL string `json:"url" validate:"required,url"`
}

type RequestLinksBatch struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}
