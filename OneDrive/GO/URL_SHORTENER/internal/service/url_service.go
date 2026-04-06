package service

type UrlRepo interface {
}
type UrlService struct {
	repo UrlRepo
}

func NewUrlService(r UrlRepo) *UrlService {
	return &UrlService{
		repo: r,
	}
}
