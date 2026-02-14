package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/james-wukong/orders-api/internal/infrastructure/security"
)

func (m *Manager) JWTAuthMiddleware(jwtManager *security.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Get the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization header is required"})
			return
		}

		// 2. Check for the "Bearer " prefix
		fields := strings.Fields(authHeader)
		if len(fields) < 2 || strings.ToLower(fields[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			return
		}

		// 3. Verify the token
		tokenString := fields[1]
		claims, err := jwtManager.VerifyToken(tokenString)
		if err != nil {
			// You can differentiate between "expired" and "invalid" if needed
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		// 4. Set the UserID in the context for downstream handlers
		// This allows handlers to do: c.GetString("user_id")
		c.Set("user_id", claims.UserID)

		// 5. Continue to the next handler
		c.Next()
	}
}
