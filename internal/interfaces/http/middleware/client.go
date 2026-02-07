// Package middleware provides middleware functions for handling HTTP requests in the application.
package middleware

import (
	"github.com/gin-gonic/gin"
)

func SetClientMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the User-Agent string from the request header
		userAgentString := c.Request.Header.Get("User-Agent")
		// parser := uaparser.NewFromSaved()
		// clientIP := c.ClientIP()

		if userAgentString == "" {
			// c.Set(constant.CtxDeviceInfo, utils.RedisJWTTokenKey(clientIP, constant.UnknownStr))
			c.Next()
			return
		}
		// client := parser.Parse(userAgentString)
		// clientDevice := client.Device.Family

		// Set device information in context
		// c.Set(constant.CtxDeviceInfo, utils.RedisJWTTokenKey(clientIP, clientDevice))

		c.Next()
	}
}
