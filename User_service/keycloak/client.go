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

func (c *Client) Login(ctx context.Context, username, password string) (*gocloak.JWT, error) {
	return c.gocloak.Login(ctx, c.clientID, c.clientSecret, c.realm, username, password)
}

func (c *Client) GetAdminToken(ctx context.Context) (string, error) {
	jwt, err := c.gocloak.LoginClient(ctx, c.clientID, c.clientSecret, c.realm)
	if err != nil {
		return "", err
	}
	return jwt.AccessToken, nil
}

func (c *Client) RegisterUser(ctx context.Context, token string, username, email, password string) error {
	user := gocloak.User{
		Username: &username,
		Email:    &email,
		Enabled:  gocloak.BoolP(true),
	}
	
	userID, err := c.gocloak.CreateUser(ctx, token, c.realm, user)
	if err != nil {
		return err
	}

	return c.gocloak.SetPassword(ctx, token, userID, c.realm, password, false)
}
