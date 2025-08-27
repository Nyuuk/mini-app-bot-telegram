package middlewares

import (
	"encoding/json"
	"time"

	"github.com/Nyuuk/mini-app-bot-telegram/backend/app/pkg/helpers"
	"github.com/gofiber/fiber/v2"
)

// LoggingMiddleware logs all HTTP requests with detailed information
func LoggingMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Capture request body
		var requestBody interface{}
		if c.Method() != "GET" && c.Method() != "HEAD" {
			body := c.Body()
			if len(body) > 0 {
				// Try to parse as JSON
				var jsonBody interface{}
				if err := json.Unmarshal(body, &jsonBody); err == nil {
					requestBody = jsonBody
				} else {
					requestBody = string(body)
				}
			}
		}

		// Process request
		err := c.Next()

		// Extract user ID from context if available
		userID := "anonymous"
		if user := c.Locals("user_id"); user != nil {
			if userStr, ok := user.(string); ok {
				userID = userStr
			} else if userUint, ok := user.(uint); ok {
				userID = string(rune(userUint))
			}
		}

		// Log the request using our structured logger
		helpers.LogRequest(
			c.Method(),
			c.Path(),
			c.IP(),
			c.Get("User-Agent"),
			time.Since(start),
			c.Response().StatusCode(),
			userID,
			requestBody,
			nil, // No response body for simplicity
		)

		return err
	}
}

// DatabaseLoggingMiddleware logs database operations
func DatabaseLoggingMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Add database logging to context
		c.Locals("db_logger", func(operation, table string, rowsAffected int64, err error) {
			helpers.LogDatabase(operation, table, time.Since(start), rowsAffected, err)
		})

		return c.Next()
	}
}

// AuthLoggingMiddleware logs authentication events
func AuthLoggingMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Log authentication attempt
		helpers.LogAuth("login_attempt", "anonymous", false, map[string]interface{}{
			"ip_address": c.IP(),
			"user_agent": c.Get("User-Agent"),
			"path":       c.Path(),
		})

		err := c.Next()

		// Log authentication result
		if userID := c.Locals("user_id"); userID != nil {
			helpers.LogAuth("login_success", userID.(string), true, map[string]interface{}{
				"ip_address": c.IP(),
				"user_agent": c.Get("User-Agent"),
			})
		}

		return err
	}
}

// PerformanceLoggingMiddleware logs performance metrics
func PerformanceLoggingMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		err := c.Next()

		duration := time.Since(start)

		// Log performance if request takes longer than 1 second
		if duration > time.Second {
			helpers.LogPerformance(
				"http_request",
				duration,
				c.Path(),
				map[string]interface{}{
					"method":     c.Method(),
					"status":     c.Response().StatusCode(),
					"ip_address": c.IP(),
				},
			)
		}

		return err
	}
}

// SecurityLoggingMiddleware logs security-related events
func SecurityLoggingMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Log potential security events
		userAgent := c.Get("User-Agent")
		if userAgent == "" {
			helpers.LogSecurity("missing_user_agent", "anonymous", c.IP(), map[string]interface{}{
				"path": c.Path(),
			})
		}

		// Log failed authentication attempts
		if c.Response().StatusCode() == 401 {
			helpers.LogSecurity("failed_auth", "anonymous", c.IP(), map[string]interface{}{
				"path":       c.Path(),
				"user_agent": userAgent,
			})
		}

		return c.Next()
	}
}
