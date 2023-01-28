package requests

type RequestLink struct {
	URL string `json:"url" validate:"required,url"`
}
