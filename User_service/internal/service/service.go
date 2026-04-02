package service

import (
	"context"
	"errors"

	"github.com/baigel/lms/user-service/keycloak"
)

type Service struct {
	kc *keycloak.Client
}

func New(kc *keycloak.Client) *Service {
	return &Service{kc: kc}
}


type TokenResponse struct {
	AccessToken      string `json:"access_token"`
	RefreshToken     string `json:"refresh_token"`
	ExpiresIn        int    `json:"expires_in"`
	RefreshExpiresIn int    `json:"refresh_expires_in"`
}

func (s *Service) Login(ctx context.Context, username, password string) (*TokenResponse, error) {
	if username == "" || password == "" {
		return nil, errors.New("username and password are required")
	}
	jwt, err := s.kc.Login(ctx, username, password)
	if err != nil {
		return nil, err
	}
	return &TokenResponse{
		AccessToken:      jwt.AccessToken,
		RefreshToken:     jwt.RefreshToken,
		ExpiresIn:        jwt.ExpiresIn,
		RefreshExpiresIn: jwt.RefreshExpiresIn,
	}, nil
}

//LT12:
func (s *Service) Refresh(ctx context.Context, refreshToken string) (*TokenResponse, error) {
	if refreshToken == "" {
		return nil, errors.New("refresh_token is required")
	}
	jwt, err := s.kc.RefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, err
	}
	return &TokenResponse{
		AccessToken:      jwt.AccessToken,
		RefreshToken:     jwt.RefreshToken,
		ExpiresIn:        jwt.ExpiresIn,
		RefreshExpiresIn: jwt.RefreshExpiresIn,
	}, nil
}

func (s *Service) Logout(ctx context.Context, refreshToken string) error {
	if refreshToken == "" {
		return errors.New("refresh_token is required")
	}
	return s.kc.Logout(ctx, refreshToken)
}

// Register — LT13: создать пользователя (вызывается только из admin-эндпоинта)
func (s *Service) Register(ctx context.Context, username, email, password string) error {
	if username == "" || email == "" || password == "" {
		return errors.New("username, email and password are required")
	}
	adminToken, err := s.kc.GetAdminToken(ctx)
	if err != nil {
		return err
	}
	return s.kc.RegisterUser(ctx, adminToken, username, email, password)
}
