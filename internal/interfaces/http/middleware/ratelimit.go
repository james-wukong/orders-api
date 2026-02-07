package middleware

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

var limiters = make(map[string]*rate.Limiter)

func RateLimiterMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		if _, exists := limiters[ip]; !exists {
			limiters[ip] = rate.NewLimiter(rate.Limit(5), 10)
		}
		limiter := limiters[ip]

		if !limiter.Allow() {
			// utils.ReturnErrorResponse(c, http.StatusTooManyRequests, "Too Many Requests", nil)
			c.Abort()
			return
		}
		c.Next()
	}
}
