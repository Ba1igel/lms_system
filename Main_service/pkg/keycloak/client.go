package keycloak

import (
	"context"

	"github.com/Nerzal/gocloak/v13"
)

type Client struct {
	gc           *gocloak.GoCloak
	clientID     string
	clientSecret string
	realm        string
}

func New(url, clientID, clientSecret, realm string) *Client {
	return &Client{
		gc:           gocloak.NewClient(url),
		clientID:     clientID,
		clientSecret: clientSecret,
		realm:        realm,
	}
}

// IntrospectToken проверяет активность токена через Keycloak introspect endpoint.
func (c *Client) IntrospectToken(ctx context.Context, accessToken string) (*gocloak.IntroSpectTokenResult, error) {
	return c.gc.RetrospectToken(ctx, accessToken, c.clientID, c.clientSecret, c.realm)
}
