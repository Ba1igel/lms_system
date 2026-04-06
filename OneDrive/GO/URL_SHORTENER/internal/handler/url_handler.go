package handler

type UrlService interface {
}
type UrlHandler struct {
	service UrlService
}

func NewUrlHandler(s UrlService) *UrlHandler {
	return &UrlHandler{
		service: s,
	}
}
