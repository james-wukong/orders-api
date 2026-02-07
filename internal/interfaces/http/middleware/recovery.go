package middleware

import (
	"fmt"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Convert the recovered value 'r' to an error type
				var recoveredErr error
				switch v := err.(type) {
				case error:
					recoveredErr = v // It's already an error
				case string:
					recoveredErr = fmt.Errorf("panic: %s", v) // Convert string to error
				default:
					recoveredErr = fmt.Errorf("panic: %v", v) // Convert any other type to error
				}
				// Log the panic
				log.Error().
					Interface("panic_value", recoveredErr). // The value passed to panic()
					Bytes("stack_trace", debug.Stack()).    // Get the full stack trace
					Str("path", c.Request.URL.Path).
					Str("method", c.Request.Method).
					Msg("Panic recovered in handler")

				// Depending on the mode, you might choose to return different responses.
				if gin.Mode() == gin.DebugMode {
					// In debug mode, send a more informative error to the client
					// utils.ReturnAutoErrorResponse(c, recoveredErr, gin.H{
					// 	"error":   "Internal Server Error - Debug Mode",
					// 	"message": fmt.Sprintf("Panic: %v", recoveredErr),
					// 	"stack":   string(debug.Stack()), // Optionally send stack in debug
					// })
				} else {
					// In release mode, send a generic error to the client
					// utils.ReturnAutoErrorResponse(c, recoveredErr, gin.H{
					// 	"error":   "Internal Server Error",
					// 	"message": "Something went wrong on our end.",
					// })
				}
			}
		}()
		// Proceed with the next handler
		c.Next()
	}
}
