package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/Nerzal/gocloak/v13"
	"github.com/baigel/lms/user-service/internal/config"
)

type contextKey string

const UserContextKey = contextKey("user")

// AuthMiddleware verifies the Keycloak JWT token
func AuthMiddleware(cfg *config.Config, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		token := parts[1]

		client := gocloak.NewClient(cfg.KeycloakURL)
		
		rptResult, err := client.RetrospectToken(r.Context(), token, cfg.KeycloakClientID, cfg.KeycloakSecret, cfg.KeycloakRealm)
		if err != nil {
			http.Error(w, "Failed to verify token", http.StatusUnauthorized)
			return
		}

		if !*rptResult.Active {
			http.Error(w, "Token is inactive or expired", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserContextKey, token)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
