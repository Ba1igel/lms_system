package middleware

import (
	"encoding/base64"
	"encoding/json"
	"strings"

	"github.com/baigel/lms/main-service/pkg/apperror"
	"github.com/baigel/lms/main-service/pkg/keycloak"
	"github.com/gin-gonic/gin"
)

// Context keys для хранения данных токена в gin.Context.
const (
	ContextKeyUserID   = "userID"
	ContextKeyUsername = "username"
	ContextKeyRoles    = "roles"
)

type jwtClaims struct {
	Sub               string      `json:"sub"`
	PreferredUsername string      `json:"preferred_username"`
	RealmAccess       realmAccess `json:"realm_access"`
}

type realmAccess struct {
	Roles []string `json:"roles"`
}

// RequireAuth проверяет наличие валидного активного токена.
// Кладёт userID, username и roles в gin.Context для использования в хендлерах.
func RequireAuth(kc *keycloak.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, ok := extractBearer(c)
		if !ok {
			_ = c.Error(apperror.Unauthorized("authorization header missing or invalid"))
			c.Abort()
			return
		}

		result, err := kc.IntrospectToken(c.Request.Context(), token)
		if err != nil || result == nil || result.Active == nil || !*result.Active {
			_ = c.Error(apperror.Unauthorized("token is invalid or expired"))
			c.Abort()
			return
		}

		claims, err := parseJWTClaims(token)
		if err != nil {
			_ = c.Error(apperror.Unauthorized("failed to parse token claims"))
			c.Abort()
			return
		}

		c.Set(ContextKeyUserID, claims.Sub)
		c.Set(ContextKeyUsername, claims.PreferredUsername)
		c.Set(ContextKeyRoles, claims.RealmAccess.Roles)

		c.Next()
	}
}

// RequireAdmin проверяет валидность токена И наличие роли "admin".
func RequireAdmin(kc *keycloak.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, ok := extractBearer(c)
		if !ok {
			_ = c.Error(apperror.Unauthorized("authorization header missing or invalid"))
			c.Abort()
			return
		}

		result, err := kc.IntrospectToken(c.Request.Context(), token)
		if err != nil || result == nil || result.Active == nil || !*result.Active {
			_ = c.Error(apperror.Unauthorized("token is invalid or expired"))
			c.Abort()
			return
		}

		claims, err := parseJWTClaims(token)
		if err != nil {
			_ = c.Error(apperror.Unauthorized("failed to parse token claims"))
			c.Abort()
			return
		}

		if !hasRole(claims.RealmAccess.Roles, "admin") {
			_ = c.Error(apperror.Forbidden("forbidden: admin role required"))
			c.Abort()
			return
		}

		c.Set(ContextKeyUserID, claims.Sub)
		c.Set(ContextKeyUsername, claims.PreferredUsername)
		c.Set(ContextKeyRoles, claims.RealmAccess.Roles)

		c.Next()
	}
}

func extractBearer(c *gin.Context) (string, bool) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return "", false
	}
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return "", false
	}
	return parts[1], true
}

// parseJWTClaims декодирует payload JWT без верификации подписи.
// Верификация активности токена уже выполнена через Keycloak introspect.
func parseJWTClaims(tokenStr string) (*jwtClaims, error) {
	parts := strings.Split(tokenStr, ".")
	if len(parts) != 3 {
		return nil, &parseError{"invalid JWT format"}
	}

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

	var claims jwtClaims
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
