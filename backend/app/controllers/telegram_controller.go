package controllers

import (
	"strconv"

	"github.com/Nyuuk/mini-app-bot-telegram/backend/app/payloads"
	"github.com/Nyuuk/mini-app-bot-telegram/backend/app/pkg/database"
	"github.com/Nyuuk/mini-app-bot-telegram/backend/app/pkg/errors"
	"github.com/Nyuuk/mini-app-bot-telegram/backend/app/pkg/helpers"
	"github.com/Nyuuk/mini-app-bot-telegram/backend/app/services"
	"github.com/gofiber/fiber/v2"
)

type TelegramController struct {
	TelegramService services.TelegramService
}

func (t *TelegramController) CreateNewUserForNowUserActive(c *fiber.Ctx) error {
	helpers.LogDebug("CreateNewUserForNowUserActive", "TelegramController: CreateNewUserForNowUserActive validate body", nil, c)
	var payload payloads.CreateNewTelegramPayload
	if err := helpers.ValidateBody(&payload, c); err != nil {
		helpers.LogError(err, "CreateNewUserForNowUserActive", "TelegramController: CreateNewUserForNowUserActive Error validating body", nil, c)
		if customErr, ok := err.(helpers.Error); ok {
			return helpers.ResponseErrorBadRequest(c, customErr.Message, customErr.Data)
		}
	}
	helpers.Logger.Debug().Msg("TelegramController: CreateNewUserForNowUserActive validate body success")
	helpers.Logger.Debug().Msg("TelegramController: CreateNewUserForNowUserActive calling service")
	tx := database.ClientPostgres.Begin()
	defer tx.Rollback()

	if err := t.TelegramService.CreateNewUserForNowUserActive(&payload, c, tx); err != nil {
		helpers.Logger.Error().Err(err).Msg("TelegramController: CreateNewUserForNowUserActive Error creating new user for now user active")
		return helpers.HandleTelegramError(c, err)
	}
	return nil
}

func (t *TelegramController) DeleteByTelegramID(c *fiber.Ctx) error {
	helpers.MyLogger("debug", "TelegramAccountLink", "DeleteByTelegramId", "controller", "start delete telegram user by telegram ID", nil, c)
	telegramID := c.Params("id")
	telegramIDInt64, err := strconv.ParseInt(telegramID, 10, 64)
	if err != nil {
		helpers.MyLogger("error", "TelegramAccountLink", "DeleteByTelegramId", "controller", "error parse telegram ID", map[string]interface{}{
			"error": err.Error(),
		}, c)
		return helpers.TelegramErrorResponse(c, errors.ErrTelegramIDParseError, "Invalid telegram ID. Please provide a valid numeric Telegram user ID.", map[string]interface{}{
			"details": "Telegram ID must be a valid integer",
		})
	}
	tx := database.ClientPostgres.Begin()
	defer tx.Rollback()

	helpers.MyLogger("debug", "TelegramAccountLink", "DeleteByTelegramId", "controller", "start calling service for delete telegram user by telegram ID", nil, c)
	if err := t.TelegramService.DeleteByTelegramID(telegramIDInt64, c, tx); err != nil {
		return helpers.HandleTelegramError(c, err)
	}
	return nil
}

// find by user id
func (t *TelegramController) FindByUserID(c *fiber.Ctx) error {
	helpers.MyLogger("debug", "TelegramAccountLink", "FindByUserId", "controller", "start find telegram user by user ID", nil, c)
	userID := helpers.GetCurrentUserID(c)
	
	if err := t.TelegramService.FindByUserID(userID, c, database.ClientPostgres); err != nil {
		helpers.MyLogger("error", "TelegramAccountLink", "FindByUserId", "controller", "error find telegram user by user ID", map[string]interface{}{
			"error": err.Error(),
			"user_id": userID,
		}, c)
		return helpers.HandleTelegramError(c, err)
	}
	return nil
}

func (t *TelegramController) FindByTelegramID(c *fiber.Ctx) error {
	helpers.MyLogger("debug", "TelegramAccountLink", "FindByTelegramId", "controller", "start find telegram user by telegram ID", nil, c)
	telegramID := c.Params("id")
	telegramIDInt64, err := strconv.ParseInt(telegramID, 10, 64)
	if err != nil {
		helpers.MyLogger("error", "TelegramAccountLink", "FindByTelegramId", "controller", "error parse telegram ID", map[string]interface{}{
			"error": err.Error(),
		}, c)
		return helpers.TelegramErrorResponse(c, errors.ErrTelegramIDParseError, "Invalid telegram ID. Please provide a valid numeric Telegram user ID.", map[string]interface{}{
			"details": "Telegram ID must be a valid integer",
		})
	}
	tx := database.ClientPostgres

	helpers.MyLogger("debug", "TelegramAccountLink", "FindByTelegramId", "controller", "start calling service for find telegram user by telegram ID", map[string]interface{}{
		"telegram_id": telegramIDInt64,
	}, c)
	if err := t.TelegramService.FindByTelegramID(telegramIDInt64, c, tx); err != nil {
		helpers.MyLogger("error", "TelegramAccountLink", "FindByTelegramId", "controller", "error find telegram user by telegram ID", map[string]interface{}{
			"error": err.Error(),
			"telegram_id": telegramIDInt64,
		}, c)
		return helpers.HandleTelegramError(c, err)
	}
	return nil
}

func (t *TelegramController) UpdateByTelegramID(c *fiber.Ctx) error {
	helpers.MyLogger("debug", "TelegramAccountLink", "UpdateByTelegramId", "controller", "start update telegram user by telegram ID", nil, c)
	telegramID := c.Params("id")
	telegramIDInt64, err := strconv.ParseInt(telegramID, 10, 64)
	if err != nil {
		helpers.MyLogger("error", "TelegramAccountLink", "UpdateByTelegramId", "controller", "error parse telegram ID", map[string]interface{}{
			"error": err.Error(),
		}, c)
		return helpers.TelegramErrorResponse(c, errors.ErrTelegramIDParseError, "Invalid telegram ID. Please provide a valid numeric Telegram user ID.", map[string]interface{}{
			"details": "Telegram ID must be a valid integer",
		})
	}
	var payload payloads.UpdateTelegramPayload
	if err := helpers.ValidateBody(&payload, c); err != nil {
		helpers.MyLogger("error", "TelegramAccountLink", "UpdateByTelegramId", "controller", "error validate body", map[string]interface{}{
			"error": err.Error(),
		}, c)
		if customErr, ok := err.(helpers.Error); ok {
			return helpers.ResponseErrorBadRequest(c, customErr.Message, customErr.Data)
		}
		return helpers.ResponseErrorBadRequest(c, "Invalid update payload. Please check your request data and try again.", nil)
	}
	tx := database.ClientPostgres.Begin()
	defer tx.Rollback()

	helpers.MyLogger("debug", "TelegramAccountLink", "UpdateByTelegramId", "controller", "start calling service for update telegram user by telegram ID", map[string]interface{}{
		"telegram_id": telegramIDInt64,
	}, c)
	if err := t.TelegramService.UpdateByTelegramID(telegramIDInt64, &payload, c, tx); err != nil {
		helpers.MyLogger("error", "TelegramAccountLink", "UpdateByTelegramId", "controller", "error update telegram user by telegram ID", map[string]interface{}{
			"error": err.Error(),
			"telegram_id": telegramIDInt64,
		}, c)
		return helpers.HandleTelegramError(c, err)
	}
	return nil
}