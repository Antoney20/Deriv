package middleware

import (
	"log"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	Reset  = "\033[0m"
	White  = "\033[37m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start)
		statusCode := c.Writer.Status()
		method := c.Request.Method
		path := c.Request.URL.Path

		userAgent := c.Request.UserAgent()
		ipAddress := c.ClientIP()

		deviceType := getDeviceType(userAgent)

		var color string
		if statusCode >= 200 && statusCode < 300 {
			color = White 
		} else if statusCode >= 400 && statusCode < 500 {
			color = Yellow 
		} else if statusCode >= 500 {
			color = Red 
		}

		logMessage := formatLogMessage(method, path, statusCode, duration, deviceType, ipAddress)
		log.Print(color + logMessage + Reset)

		// log.Printf("Request: %s %s | Status: %d | Duration: %v | Device: %s | IP: %s\n",
		// 	method, path, statusCode, duration, deviceType, ipAddress)
	}
}

func getDeviceType(userAgent string) string {
	if strings.Contains(userAgent, "Android") {
		return "Android"
	} else if strings.Contains(userAgent, "iPhone") {
		return "iPhone"
	} else if strings.Contains(userAgent, "Windows") {
		return "Windows"
	} else if strings.Contains(userAgent, "Macintosh") {
		return "Mac"
	} else if strings.Contains(userAgent, "Linux") {
		return "Linux"
	}
	return "Unknown"
}

func formatLogMessage(method, path string, statusCode int, duration time.Duration, deviceType, ipAddress string) string {
	return "Request: " + method + " " + path + " | Status: " + string(statusCode) +
		" | Duration: " + duration.String() + " | Device: " + deviceType + " | IP: " + ipAddress + "\n"
}
