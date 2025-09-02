package helpers

import (
	"fmt"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type Payload interface {
	CustomErrorsMessage(validator.ValidationErrors) []map[string]string
}

func ValidateBody(payload Payload, c *fiber.Ctx) error {
	// log.Debug("Validating request body ", payload)
	// Logger.Debug().Interface("payload", payload).Msg("Validating request body")
	validate := validator.New()

	// Parse the request body
	if err := c.BodyParser(&payload); err != nil {
		// log.Error("Error parsing re quest body: ", err)
		// Logger.Error().Err(err).Msg("Error parsing request body")
		// ResponseErrorBadRequest(c, "Invalid payload", err)
		return err
	}
	// log.Debug("Validating after body parser ", payload)
	// Logger.Debug().Interface("payload", payload).Msg("Validating after body parser")

	// Validate the user struct
	if err := validate.Struct(payload); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		data := payload.CustomErrorsMessage(validationErrors)
		// log.Debug("Validation errors ", data)
		// Logger.Debug().Interface("validation_errors", data).Msg("Validation errors")
		// log.Debug("Validation errors ", validationErrors)
		return ErrorClient("Invalid payload ", fiber.StatusBadRequest, data)
		// return ErrorClient("Invalid payload ", fiber.StatusBadRequest, validationErrors)
	}

	return nil
}

func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// GetTimezone returns the configured timezone or defaults to Asia/Jakarta
func GetTimezone() *time.Location {
	timezone := GetEnv("TIMEZONE", "Asia/Jakarta")
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		// Fallback to Asia/Jakarta if invalid timezone
		loc, _ = time.LoadLocation("Asia/Jakarta")
	}
	return loc
}

// ParseDateWithTimezone parses date string with configured timezone
func ParseDateWithTimezone(dateStr string) (time.Time, error) {
	loc := GetTimezone()
	return time.ParseInLocation("2006-01-02", dateStr, loc)
}

// ParseDateTimeWithTimezone parses datetime string with configured timezone
// Supports multiple formats: "2006-01-02T15:04:05", "2006-01-02 15:04:05", "15:04:05", "15:04"
func ParseDateTimeWithTimezone(dateTimeStr string) (time.Time, error) {
	loc := GetTimezone()

	// Try different datetime formats
	formats := []string{
		"2006-01-02T15:04:05",       // ISO format without timezone
		"2006-01-02 15:04:05",       // Space separated format
		"2006-01-02T15:04:05Z07:00", // With timezone (fallback)
	}

	for _, format := range formats {
		if t, err := time.ParseInLocation(format, dateTimeStr, loc); err == nil {
			return t, nil
		}
	}

	// Try time-only format (HH:MM:SS) - use base date (1900-01-01) for time-only storage
	if t, err := time.ParseInLocation("15:04:05", dateTimeStr, loc); err == nil {
		// Use base date 1900-01-01 with the provided time for time-only storage
		return time.Date(1900, 1, 1, t.Hour(), t.Minute(), t.Second(), 0, loc), nil
	}

	// Try time-only format (HH:MM) - use base date (1900-01-01) for time-only storage
	if t, err := time.ParseInLocation("15:04", dateTimeStr, loc); err == nil {
		// Use base date 1900-01-01 with the provided time for time-only storage
		return time.Date(1900, 1, 1, t.Hour(), t.Minute(), 0, 0, loc), nil
	}

	// If all formats fail, return error
	return time.Time{}, fmt.Errorf("invalid datetime format: %s", dateTimeStr)
}

// ParseTimeWithTimezone parses time string with configured timezone
// Returns time with base date (1900-01-01) for time-only storage
func ParseTimeWithTimezone(timeStr string) (time.Time, error) {
	loc := GetTimezone()

	// Try HH:MM:SS format
	if t, err := time.ParseInLocation("15:04:05", timeStr, loc); err == nil {
		return time.Date(1900, 1, 1, t.Hour(), t.Minute(), t.Second(), 0, loc), nil
	}

	// Try HH:MM format
	if t, err := time.ParseInLocation("15:04", timeStr, loc); err == nil {
		return time.Date(1900, 1, 1, t.Hour(), t.Minute(), 0, 0, loc), nil
	}

	return time.Time{}, fmt.Errorf("invalid time format: %s", timeStr)
}

// IsTimeOnlyFormat checks if the string is time-only format (HH:MM or HH:MM:SS)
func IsTimeOnlyFormat(timeStr string) bool {
	loc := GetTimezone()

	// Check if it's HH:MM:SS format
	if _, err := time.ParseInLocation("15:04:05", timeStr, loc); err == nil {
		return true
	}

	// Check if it's HH:MM format
	if _, err := time.ParseInLocation("15:04", timeStr, loc); err == nil {
		return true
	}

	return false
}

// NowWithTimezone returns current time in configured timezone
func NowWithTimezone() time.Time {
	loc := GetTimezone()
	return time.Now().In(loc)
}

// ConvertToTimezone converts time to configured timezone
func ConvertToTimezone(t time.Time) time.Time {
	loc := GetTimezone()
	return t.In(loc)
}
