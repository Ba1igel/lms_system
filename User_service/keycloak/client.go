package keycloak

import (
	"context"

	"github.com/Nerzal/gocloak/v13"
)

type Client struct {
	gocloak      *gocloak.GoCloak
	clientID     string
	clientSecret string
	realm        string
}

func New(url, clientID, clientSecret, realm string) *Client {
	return &Client{
		gocloak:      gocloak.NewClient(url),
		clientID:     clientID,
		clientSecret: clientSecret,
		realm:        realm,
	}
}

// Login — получить access + refresh токены по логину/паролю (LT11)
func (c *Client) Login(ctx context.Context, username, password string) (*gocloak.JWT, error) {
	return c.gocloak.Login(ctx, c.clientID, c.clientSecret, c.realm, username, password)
}

//  (LT12)
func (c *Client) RefreshToken(ctx context.Context, refreshToken string) (*gocloak.JWT, error) {
	return c.gocloak.RefreshToken(ctx, refreshToken, c.clientID, c.clientSecret, c.realm)
}

// отозвать refresh_token (LT12)
func (c *Client) Logout(ctx context.Context, refreshToken string) error {
	return c.gocloak.Logout(ctx, c.clientID, c.clientSecret, c.realm, refreshToken)
}

// ervice-account токен клиента (LT13)
func (c *Client) GetAdminToken(ctx context.Context) (string, error) {
	jwt, err := c.gocloak.LoginClient(ctx, c.clientID, c.clientSecret, c.realm)
	if err != nil {
		return "", err
	}
	return jwt.AccessToken, nil
}

// RegisterUser — создать пользователя в Keycloak (LT13)
func (c *Client) RegisterUser(ctx context.Context, adminToken, username, email, password string) error {
	enabled := true
	user := gocloak.User{
		Username: &username,
		Email:    &email,
		Enabled:  &enabled,
	}
	userID, err := c.gocloak.CreateUser(ctx, adminToken, c.realm, user)
	if err != nil {
		return err
	}
	return c.gocloak.SetPassword(ctx, adminToken, userID, c.realm, password, false)
}

//чекаем активность токена  активность токена (middleware, LT13)
func (c *Client) IntrospectToken(ctx context.Context, accessToken string) (*gocloak.IntroSpectTokenResult, error) {
	return c.gocloak.RetrospectToken(ctx, accessToken, c.clientID, c.clientSecret, c.realm)
}
