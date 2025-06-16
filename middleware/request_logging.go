package middleware

import (
	"fmt"
	"intern-project-v2/logger"
	"time"

	"github.com/gin-gonic/gin"
)

func RequestLogging() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		c.Next()

		duration := time.Since(startTime)
		statusCode := c.Writer.Status()
		responseSize := c.Writer.Size()
		durationMs := formatDuration(duration)
		responseSizeFormatted := formatBytes(responseSize)

		switch {
		case statusCode >= 500:
			logger.Error("Request ended with error", "Time:", durationMs, "Method:", c.Request.Method, "Path:", c.Request.URL.Path, "From:", c.ClientIP(), "User-Agent:", c.Request.UserAgent(), "Status Code:", statusCode, "Response Size:", responseSizeFormatted)
		case statusCode >= 400:
			logger.Warn("Request ended with warning", "Time:", durationMs, "Method:", c.Request.Method, "Path:", c.Request.URL.Path, "From:", c.ClientIP(), "User-Agent:", c.Request.UserAgent(), "Status Code:", statusCode, "Response Size:", responseSizeFormatted)
		default:
			logger.Info("Request ended successfully", "Time:", durationMs, "Method:", c.Request.Method, "Path:", c.Request.URL.Path, "From:", c.ClientIP(), "User-Agent:", c.Request.UserAgent(), "Status Code:", statusCode, "Response Size:", responseSizeFormatted)
		}

	}
}

func formatDuration(duration time.Duration) string {
	if duration < time.Microsecond {
		return fmt.Sprintf("%.2fns", float64(duration.Nanoseconds()))
	} else if duration < time.Millisecond {
		return fmt.Sprintf("%.2fμs", float64(duration.Nanoseconds())/1000)
	} else if duration < time.Second {
		return fmt.Sprintf("%.2fms", float64(duration.Nanoseconds())/1000000)
	} else {
		return fmt.Sprintf("%.2fs", duration.Seconds())
	}
}

func formatBytes(bytes int) string {
	if bytes < 1024 {
		return fmt.Sprintf("%dB", bytes)
	} else if bytes < 1024*1024 {
		return fmt.Sprintf("%.2fKB", float64(bytes)/1024)
	} else {
		return fmt.Sprintf("%.2fMB", float64(bytes)/(1024*1024))
	}
}
