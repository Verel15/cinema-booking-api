package middleware

import (
	"cinema-booking-api/pkg/logger"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// Generate Request ID
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		c.Header("X-Request-ID", requestID)
		c.Set("request_id", requestID)

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)

		// Extract user_id if available (placeholder for now)
		userID, _ := c.Get("user_id")
		if userID == nil {
			userID = ""
		}

		// Log fields
		fields := []zap.Field{
			zap.String("request_id", requestID),
			zap.String("user_id", userID.(string)),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.Int("status", c.Writer.Status()),
			zap.String("latency", latency.String()),
		}

		// Include query if present
		if query != "" {
			fields = append(fields, zap.String("query", query))
		}

		// Check for errors in Gin context
		if len(c.Errors) > 0 {
			for _, e := range c.Errors {
				logger.Log.Error(e.Error(), fields...)
			}
		} else {
			// Normal request logging
			if c.Writer.Status() >= 400 {
				logger.Log.Warn("request failed", fields...)
			} else {
				logger.Log.Info("request processed", fields...)
			}
		}
	}
}
