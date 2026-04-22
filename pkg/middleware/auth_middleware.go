package middleware

import (
	"cinema-booking-api/internal/user/domain"
	"cinema-booking-api/pkg/response"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	validateToken func(token string) (*domain.User, error)
}

func NewAuthMiddleware(validateToken func(token string) (*domain.User, error)) *AuthMiddleware {
	return &AuthMiddleware{validateToken: validateToken}
}

func (m *AuthMiddleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, http.StatusUnauthorized, "authorization header is required")
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Error(c, http.StatusUnauthorized, "invalid authorization header format")
			c.Abort()
			return
		}

		token := parts[1]

		// Validate token
		user, err := m.validateToken(token)
		if err != nil {
			response.Error(c, http.StatusUnauthorized, "invalid token")
			c.Abort()
			return
		}

		// Set user in context as map for easy access by RBAC
		userMap := map[string]interface{}{
			"id":       user.ID,
			"email":    user.Email,
			"role":     string(user.Role),
			"provider": user.Provider,
		}
		c.Set("user", userMap)
		c.Next()
	}
}
