package logger

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RequestIDKey is the key used to store the request ID in the context
const RequestIDKey = "request_id"

// RequestLogger is a middleware that logs HTTP requests
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Generate a unique request ID
		requestID := uuid.New().String()
		c.Set(RequestIDKey, requestID)
		c.Header("X-Request-ID", requestID)

		// Start timer
		start := time.Now()

		// Process request
		c.Next()

		// Calculate request processing time
		duration := time.Since(start)

		// Get status code
		statusCode := c.Writer.Status()

		// Log request details
		fmt.Printf("[%s] %s %s %d %dms %s %s\n",
			requestID,
			c.Request.Method,
			c.Request.URL.Path,
			statusCode,
			duration.Milliseconds(),
			c.ClientIP(),
			c.Request.UserAgent(),
		)
	}
}

// GetRequestID returns the request ID from the context
func GetRequestID(c *gin.Context) string {
	if id, exists := c.Get(RequestIDKey); exists {
		return id.(string)
	}
	return ""
}
