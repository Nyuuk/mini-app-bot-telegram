package helpers

import (
	"testing"

	telegram_errors "github.com/Nyuuk/mini-app-bot-telegram/backend/app/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestTelegramErrorCreation(t *testing.T) {
	// Test Telegram error creation
	telegramErr := telegram_errors.NewTelegramError(
		telegram_errors.ErrTelegramUserNotFound,
		"Telegram user not found",
		map[string]interface{}{"telegram_id": int64(123456789)},
	)

	assert.Equal(t, telegram_errors.ErrTelegramUserNotFound, telegramErr.Code)
	assert.Equal(t, "Telegram user not found", telegramErr.Message)
	assert.Equal(t, map[string]interface{}{"telegram_id": int64(123456789)}, telegramErr.Details)
}

func TestTelegramErrorResponseCreation(t *testing.T) {
	// Test that we can create a Telegram error response
	telegramErr := telegram_errors.NewTelegramError(
		telegram_errors.ErrTelegramInvalidID,
		"Invalid Telegram ID",
		map[string]interface{}{
			"details": "Telegram ID must be a valid integer",
		},
	)

	assert.Equal(t, telegram_errors.ErrTelegramInvalidID, telegramErr.Code)
	assert.Equal(t, "Invalid Telegram ID", telegramErr.Message)
}