package helpers

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Nyuuk/mini-app-bot-telegram/backend/app/entities"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var Logger zerolog.Logger

// InitLogger initializes the logger with proper configuration
func InitLogger() {
	// Set time format
	zerolog.TimeFieldFormat = time.RFC3339

	// Configure log level based on environment
	logLevel := getLogLevel()

	// Configure output format based on environment
	if os.Getenv("ENV") == "production" {
		// Production: JSON format for better parsing
		Logger = zerolog.New(os.Stdout).
			Level(logLevel).
			With().
			Caller().
			Timestamp().
			Logger()
	} else {
		// Development: Pretty console output
		// Logger = zerolog.New(zerolog.ConsoleWriter{
		// 	Out:        os.Stdout,
		// 	TimeFormat: time.RFC3339,
		// 	FormatLevel: func(i interface{}) string {
		// 		if ll, ok := i.(string); ok {
		// 			switch ll {
		// 			case "trace":
		// 				return "\x1b[37mTRC\x1b[0m"
		// 			case "debug":
		// 				return "\x1b[36mDBG\x1b[0m"
		// 			case "info":
		// 				return "\x1b[32mINF\x1b[0m"
		// 			case "warn":
		// 				return "\x1b[33mWRN\x1b[0m"
		// 			case "error":
		// 				return "\x1b[31mERR\x1b[0m"
		// 			case "fatal":
		// 				return "\x1b[35mFTL\x1b[0m"
		// 			case "panic":
		// 				return "\x1b[35mPNC\x1b[0m"
		// 			default:
		// 				return ll
		// 			}
		// 		}
		// 		return "???"
		// 	},
		// 	FormatMessage: func(i interface{}) string {
		// 		return "\x1b[1m" + i.(string) + "\x1b[0m"
		// 	},
		// 	FormatFieldName: func(i interface{}) string {
		// 		return "\x1b[36m" + i.(string) + "\x1b[0m="
		// 	},
		// 	FormatFieldValue: func(i interface{}) string {
		// 		if str, ok := i.(string); ok {
		// 			return "\x1b[32m" + str + "\x1b[0m "
		// 		}
		// 		// Handle non-string values safely
		// 		return "\x1b[32m" + fmt.Sprintf("%v", i) + "\x1b[0m "
		// 	},
		// }).
		Logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}).
			Level(logLevel).
			With().
			Caller().
			Timestamp().
			Logger()
	}

	// Set global logger
	log.Logger = Logger
}

// getLogLevel returns the appropriate log level based on environment
func getLogLevel() zerolog.Level {
	level := os.Getenv("LOG_LEVEL")
	switch level {
	case "trace":
		return zerolog.TraceLevel
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	case "fatal":
		return zerolog.FatalLevel
	case "panic":
		return zerolog.PanicLevel
	default:
		return zerolog.InfoLevel
	}
}

// LogRequest logs incoming HTTP requests with detailed information
func LogRequest(method, path, remoteAddr, userAgent string, duration time.Duration, statusCode int, userID string, requestBody, responseBody interface{}) {
	l := Logger.With().CallerWithSkipFrameCount(3).Logger()
	logEvent := l.Debug().
		Str("type", "http_request").
		Str("method", method).
		Str("path", path).
		Str("remote_addr", remoteAddr).
		Str("user_agent", userAgent).
		Dur("duration", duration).
		Int("status_code", statusCode).
		Str("user_id", userID)

	// Add request body if provided
	if requestBody != nil {
		body := make(map[string]interface{})
		for key, value := range requestBody.(map[string]interface{}) {
			// fmt.Println("looping request body", key, value)
			if isSensitiveField(key) {
				body[key] = "*****"
				continue
			}
			body[key] = SafeLogValue(value)
			// logEvent = logEvent.Interface(key, SafeLogValue(value))
		}
		logEvent = logEvent.Interface("request_body", body)
	}

	// Add response body if provided
	if responseBody != nil {
		logEvent = logEvent.Interface("response_body", responseBody)
	}

	logEvent.Msg("HTTP Request")
}

// LogDatabase logs database operations
func LogDatabase(operation, table string, duration time.Duration, rowsAffected int64, error error) {
	l := Logger.With().CallerWithSkipFrameCount(3).Logger()
	event := l.Info().
		Str("type", "database").
		Str("operation", operation).
		Str("table", table).
		Dur("duration", duration).
		Int64("rows_affected", rowsAffected)

	if error != nil {
		event.Err(error).Msg("Database Error")
	} else {
		event.Msg("Database Operation")
	}
}

// LogAuth logs authentication events
func LogAuth(event string, userID string, success bool, details map[string]interface{}) {
	l := Logger.With().CallerWithSkipFrameCount(4).Logger()
	logEvent := l.Debug().
		Str("type", "auth").
		Str("event", event).
		Str("user_id", userID).
		Bool("success", success)

	// Add additional details safely
	for key, value := range details {
		// Skip sensitive data
		if isSensitiveField(key) {
			continue
		}
		logEvent = logEvent.Interface(key, value)
	}

	if success {
		logEvent.Msg("Authentication Success")
	} else {
		logEvent.Msg("Authentication Failed")
	}
}

// isSensitiveField checks if a field name contains sensitive data
func isSensitiveField(fieldName string) bool {
	sensitiveFields := []string{
		"password", "token", "secret", "key", "credential",
		"api_key", "jwt_token", "access_token", "refresh_token",
	}

	fieldNameLower := strings.ToLower(fieldName)
	for _, sensitive := range sensitiveFields {
		if strings.Contains(fieldNameLower, sensitive) {
			return true
		}
	}
	return false
}

// LogBusiness logs business logic events
func LogBusiness(event string, userID string, details map[string]interface{}) {
	l := Logger.With().CallerWithSkipFrameCount(3).Logger()
	logEvent := l.Info().
		Str("type", "business").
		Str("event", event).
		Str("user_id", userID)

	// Add additional details
	for key, value := range details {
		logEvent = logEvent.Interface(key, value)
	}

	logEvent.Msg("Business Event")
}

// LogSecurity logs security-related events
func LogSecurity(event string, userID string, ipAddress string, details map[string]interface{}) {
	l := Logger.With().CallerWithSkipFrameCount(3).Logger()
	logEvent := l.Warn().
		Str("type", "security").
		Str("event", event).
		Str("user_id", userID).
		Str("ip_address", ipAddress)

	// Add additional details
	for key, value := range details {
		logEvent = logEvent.Interface(key, value)
	}

	logEvent.Msg("Security Event")
}

// LogPerformance logs performance metrics
func LogPerformance(operation string, duration time.Duration, resource string, details map[string]interface{}) {
	l := Logger.With().CallerWithSkipFrameCount(3).Logger()
	logEvent := l.Info().
		Str("type", "performance").
		Str("operation", operation).
		Dur("duration", duration).
		Str("resource", resource)

	// Add additional details
	for key, value := range details {
		logEvent = logEvent.Interface(key, value)
	}

	logEvent.Msg("Performance Metric")
}

// LogTelegram logs Telegram-specific events
func LogTelegram(event string, chatID int64, userID int64, details map[string]interface{}) {
	l := Logger.With().CallerWithSkipFrameCount(3).Logger()
	logEvent := l.Info().
		Str("type", "telegram").
		Str("event", event).
		Int64("chat_id", chatID).
		Int64("user_id", userID)

	// Add additional details
	for key, value := range details {
		logEvent = logEvent.Interface(key, value)
	}

	logEvent.Msg("Telegram Event")
}

// LogAPI logs API-related events
func LogAPI(operation string, endpoint string, method string, statusCode int, duration time.Duration, details map[string]interface{}) {
	l := Logger.With().CallerWithSkipFrameCount(3).Logger()
	logEvent := l.Info().
		Str("type", "api").
		Str("operation", operation).
		Str("endpoint", endpoint).
		Str("method", method).
		Int("status_code", statusCode).
		Dur("duration", duration)

	// Add additional details
	for key, value := range details {
		logEvent = logEvent.Interface(key, value)
	}

	logEvent.Msg("API Operation")
}

func LogInfo(event string, message string, details map[string]interface{}, c *fiber.Ctx) {
	var user_id uint
	var user_name string
	if user, ok := c.Locals("user").(entities.User); ok {
		user_id = user.ID
		user_name = user.Username
	} else {
		user_id = 0
		user_name = "anonymous"
	}
	l := Logger.With().CallerWithSkipFrameCount(3).Logger()
	logEvent := l.Info().
		Str("type", "info").
		Str("even", event).
		Uint("user_id", user_id).
		Str("user_name", user_name)

	// Add additional details
	for key, value := range details {
		logEvent = logEvent.Interface(key, SafeLogValue(value))
	}

	logEvent.Msg(message)
}

func LogDebug(event string, message string, details map[string]interface{}, c *fiber.Ctx) {
	// user_id := c.Locals("user_id").(uint)
	var user_id uint
	var user_name string
	if user, ok := c.Locals("user").(entities.User); ok {
		user_id = user.ID
		user_name = user.Username
	} else {
		user_id = 0
		user_name = "anonymous"
	}
	l := Logger.With().CallerWithSkipFrameCount(3).Logger()
	logEvent := l.Debug().
		Str("type", "debug").
		Str("event", event).
		Uint("user_id", user_id).
		Str("user_name", user_name)

		// Add additional details
	for key, value := range details {
		logEvent = logEvent.Interface(key, SafeLogValue(value))
	}
	logEvent.Msg(message)
}

// LogError logs errors with context
func LogError(err error, event string, message string, details map[string]interface{}, c *fiber.Ctx) {
	var user_id uint
	var user_name string
	if user, ok := c.Locals("user").(entities.User); ok {
		user_id = user.ID
		user_name = user.Username
	} else {
		user_id = 0
		user_name = "anonymous"
	}
	l := Logger.With().CallerWithSkipFrameCount(3).Logger()
	logEvent := l.Error().
		Err(err).
		Str("type", "error").
		Str("even", event).
		Uint("user_id", user_id).
		Str("user_name", user_name)

	// Add additional details
	for key, value := range details {
		logEvent = logEvent.Interface(key, SafeLogValue(value))
	}

	logEvent.Msg(message)
}

// SafeLogValue converts values to safe logging format
func SafeLogValue(value interface{}) interface{} {
	if value == nil {
		return "null"
	}

	switch v := value.(type) {
	case *time.Time:
		if v == nil {
			return "null"
		}
		return v.Format("2006-01-02 15:04:05")
	case bool:
		return fmt.Sprintf("%t", v)
	case string:
		if v == "" {
			return "empty"
		}
		return v
	default:
		return v
	}
}
