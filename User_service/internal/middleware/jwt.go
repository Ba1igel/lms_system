package middleware

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/baigel/lms/user-service/keycloak"
	"github.com/gin-gonic/gin"
)

type Middleware struct {
	kc *keycloak.Client
}

func New(kc *keycloak.Client) *Middleware {
	return &Middleware{kc: kc}
}

type realmAccess struct {
	Roles []string `json:"roles"`
}

type keycloakClaims struct {
	RealmAccess realmAccess `json:"realm_access"`
}

// RequireAdmin — LT13: middleware, пропускает только запросы с ролью "admin"
//
// Алгоритм:
//  1. Извлечь Bearer токен из заголовка Authorization
//  2. Проверить активность токена через Keycloak introspect endpoint
//  3. Распарсить JWT payload → получить realm_access.roles
//  4. Убедиться, что в ролях есть "admin"
func (m *Middleware) RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken, ok := extractBearerToken(c.Request)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header missing or invalid"})
			c.Abort()
			return
		}

		result, err := m.kc.IntrospectToken(c.Request.Context(), accessToken)
		if err != nil || result == nil || result.Active == nil || !*result.Active {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token is invalid or expired"})
			c.Abort()
			return
		}

		claims, err := parseJWTClaims(accessToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "failed to parse token claims"})
			c.Abort()
			return
		}

		if !hasRole(claims.RealmAccess.Roles, "admin") {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: admin role required"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// extractBearerToken — вытащить токен из "Authorization: Bearer <token>"
func extractBearerToken(r *http.Request) (string, bool) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", false
	}
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return "", false
	}
	return parts[1], true
}

// parseJWTClaims — декодирует payload JWT без верификации подписи
// JWT структура: <header>.<payload>.<signature>  (base64url encoded)
func parseJWTClaims(tokenStr string) (*keycloakClaims, error) {
	parts := strings.Split(tokenStr, ".")
	if len(parts) != 3 {
		return nil, &parseError{"invalid JWT format"}
	}

	// base64url → base64 standard (добавляем паддинг)
	payload := parts[1]
	if pad := len(payload) % 4; pad != 0 {
		payload += strings.Repeat("=", 4-pad)
	}
	payload = strings.ReplaceAll(payload, "-", "+")
	payload = strings.ReplaceAll(payload, "_", "/")

	decoded, err := base64.StdEncoding.DecodeString(payload)
	if err != nil {
		return nil, err
	}

	var claims keycloakClaims
	if err := json.Unmarshal(decoded, &claims); err != nil {
		return nil, err
	}
	return &claims, nil
}

func hasRole(roles []string, target string) bool {
	for _, r := range roles {
		if r == target {
			return true
		}
	}
	return false
}

type parseError struct{ msg string }

func (e *parseError) Error() string { return e.msg }
