package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// mlogger func
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Process the request
		c.Next()

		// details
		duration := time.Since(start)
		statusCode := c.Writer.Status()
		method := c.Request.Method
		path := c.Request.URL.Path

		log.Printf("Request: %s %s | Status: %d | Duration: %v\n", method, path, statusCode, duration)
	}
}
