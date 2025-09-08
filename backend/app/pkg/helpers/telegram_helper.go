package helpers

import (
	"github.com/Nyuuk/mini-app-bot-telegram/backend/app/pkg/errors"
	"github.com/gofiber/fiber/v2"
)

// HandleTelegramError processes Telegram-specific errors and returns appropriate HTTP responses
func HandleTelegramError(c *fiber.Ctx, err error) error {
	// Check if it's already a Telegram error
	if telegramErr, ok := err.(errors.TelegramError); ok {
		switch telegramErr.Code {
		case errors.ErrTelegramUserNotFound:
			return Response(c, fiber.StatusNotFound, telegramErr.Message, map[string]interface{}{
				"error_code": telegramErr.Code,
				"details":    telegramErr.Details,
			})
		case errors.ErrTelegramAlreadyLinked:
			return Response(c, fiber.StatusConflict, telegramErr.Message, map[string]interface{}{
				"error_code": telegramErr.Code,
				"details":    telegramErr.Details,
			})
		case errors.ErrTelegramInvalidID, errors.ErrTelegramIDParseError:
			return Response(c, fiber.StatusBadRequest, telegramErr.Message, map[string]interface{}{
				"error_code": telegramErr.Code,
				"details":    telegramErr.Details,
			})
		default:
			return Response(c, fiber.StatusBadRequest, telegramErr.Message, map[string]interface{}{
				"error_code": telegramErr.Code,
				"details":    telegramErr.Details,
			})
		}
	}
	
	// For generic errors, return internal server error
	return ResponseErrorInternal(c, err)
}

// TelegramResponse creates a standardized response for Telegram operations
func TelegramResponse(c *fiber.Ctx, status int, message string, data interface{}) error {
	return Response(c, status, message, data)
}

// TelegramErrorResponse creates a standardized error response for Telegram operations
func TelegramErrorResponse(c *fiber.Ctx, errorCode, message string, details interface{}) error {
	return Response(c, fiber.StatusBadRequest, message, map[string]interface{}{
		"error_code": errorCode,
		"details":    details,
	})
}