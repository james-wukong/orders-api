// Package middleware provides middleware functions for handling HTTP requests in the application.
package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/james-wukong/orders-api/internal/pkg/utils"
)

func (m *Manager) SetClientMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Initialize the metadata extractor and get metadata from the request
		meta := utils.NewMetadataExtractor().GetMetadata(c)

		c.Set("client_ip", meta.IP)
		c.Set("client_os", meta.OS)
		c.Set("client_device", meta.Device)
		c.Set("client_user_agent", meta.UserAgent)
		for k, v := range meta.Location {
			c.Set("client_location_"+k, v)
		}
		m.log.Info().Msgf("Client Metadata: IP=%s, OS=%s, Device=%s, UserAgent=%s, Location=%v",
			meta.IP, meta.OS, meta.Device, meta.UserAgent, meta.Location)
		c.Next()
	}
}
