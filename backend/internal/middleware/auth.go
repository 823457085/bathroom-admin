package middleware

import (
	"net/http"
	"strings"

	"bathroom-admin/pkg/jwt"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	jwtManager *jwt.JWTManager
}

func NewAuthMiddleware(jm *jwt.JWTManager) *AuthMiddleware {
	return &AuthMiddleware{jwtManager: jm}
}

func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization format"})
			return
		}

		claims, err := m.jwtManager.Validate(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}

		c.Set("user_id", claims.UserID)
		c.Next()
	}
}
