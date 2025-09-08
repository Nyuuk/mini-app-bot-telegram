package helpers

import (
	"github.com/Nyuuk/mini-app-bot-telegram/backend/app/pkg/errors"
	"github.com/gofiber/fiber/v2"
)

func Response(c *fiber.Ctx, status int, message string, data interface{}) error {
	response := fiber.Map{
		"code":    status,
		"message": message,
	}
	
	// If data is provided, include it in the response
	if data != nil {
		// Check if data contains error details
		if details, ok := data.(map[string]interface{}); ok {
			if errorCode, exists := details["error_code"]; exists {
				response["error_code"] = errorCode
			}
			if errorDetails, exists := details["details"]; exists {
				response["details"] = errorDetails
			}
			// If it's just details without error_code, include the whole data
			if _, hasErrorCode := details["error_code"]; !hasErrorCode {
				response["data"] = data
			}
		} else {
			response["data"] = data
		}
	}
	
	return c.Status(status).JSON(response)
}

func ResponseErrorInternal(c *fiber.Ctx, err any) error {
	// log.Error("Internal server error: ", err)
	// Check if it's a Telegram error
	if telegramErr, ok := err.(errors.TelegramError); ok {
		return Response(c, fiber.StatusInternalServerError, "Internal server error", map[string]interface{}{
			"error_code": telegramErr.Code,
			"details":    telegramErr.Details,
		})
	}
	return Response(c, fiber.StatusInternalServerError, "Internal server error", nil)
}

func ResponseErrorBadRequest(c *fiber.Ctx, message string, err any) error {
	// log.Error("Bad request: ", err)
	// Check if it's a Telegram error
	if telegramErr, ok := err.(errors.TelegramError); ok {
		return Response(c, fiber.StatusBadRequest, message, map[string]interface{}{
			"error_code": telegramErr.Code,
			"details":    telegramErr.Details,
		})
	}
	return Response(c, fiber.StatusBadRequest, message, err)
}

func ResponseErrorNotFound(c *fiber.Ctx, err any) error {
	// log.Error("Not found: ", err)
	// Check if it's a Telegram error
	if telegramErr, ok := err.(errors.TelegramError); ok {
		return Response(c, fiber.StatusNotFound, "Not found", map[string]interface{}{
			"error_code": telegramErr.Code,
			"details":    telegramErr.Details,
		})
	}
	return Response(c, fiber.StatusNotFound, "Not found", nil)
}
