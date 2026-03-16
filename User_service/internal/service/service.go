package service

import (
	"context"
	"github.com/baigel/lms/user-service/internal/config"
)

type Service struct {
	cfg *config.Config
}

func (s *Service) Login(context context.Context, username string, password string) (any, any) {
	panic("unimplemented")
}
func New(cfg *config.Config) *Service {
	return &Service{
		cfg: cfg,
	}
}
