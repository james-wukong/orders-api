package middleware

import "github.com/gin-gonic/gin"

func (m *Manager) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Logic using m.db or m.log
		c.Next()
	}
}
