package logger

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// HTTPMiddleware creates a Fiber middleware for logging HTTP requests
func HTTPMiddleware(logger *Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Get request ID if exists
		requestID := c.Get("X-Request-ID", "")
		if requestID == "" {
			requestID = c.Locals("requestid").(string)
		}

		// Process request
		err := c.Next()

		// Calculate duration
		duration := time.Since(start)

		// Log request
		fields := []zap.Field{
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.Int("status", c.Response().StatusCode()),
			zap.Duration("duration", duration),
			zap.String("ip", c.IP()),
			zap.String("user_agent", c.Get("User-Agent")),
		}

		if requestID != "" {
			fields = append(fields, zap.String("request_id", requestID))
		}

		// Add user_id if authenticated
		if userID := c.Locals("userID"); userID != nil {
			fields = append(fields, zap.String("user_id", userID.(string)))
		}

		// Log with appropriate level based on status code
		statusCode := c.Response().StatusCode()
		if err != nil {
			fields = append(fields, zap.Error(err))
			logger.Error("HTTP request error", fields...)
		} else if statusCode >= 500 {
			logger.Error("HTTP request completed with server error", fields...)
		} else if statusCode >= 400 {
			logger.Warn("HTTP request completed with client error", fields...)
		} else {
			logger.Info("HTTP request completed", fields...)
		}

		return err
	}
}

// RequestIDMiddleware adds a request ID to each request
func RequestIDMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		requestID := c.Get("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
		}
		c.Locals("requestid", requestID)
		c.Set("X-Request-ID", requestID)
		return c.Next()
	}
}

// generateRequestID generates a simple request ID
func generateRequestID() string {
	return time.Now().Format("20060102150405.000000")
}
