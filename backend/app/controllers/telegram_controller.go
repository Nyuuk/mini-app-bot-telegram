package controllers

import (
	"strconv"

	"github.com/Nyuuk/mini-app-bot-telegram/backend/app/payloads"
	"github.com/Nyuuk/mini-app-bot-telegram/backend/app/pkg/database"
	"github.com/Nyuuk/mini-app-bot-telegram/backend/app/pkg/helpers"
	"github.com/Nyuuk/mini-app-bot-telegram/backend/app/services"
	"github.com/gofiber/fiber/v2"
)

type TelegramController struct {
	TelegramService services.TelegramService
}

func (t *TelegramController) CreateNewUserForNowUserActive(c *fiber.Ctx) error {
	// helpers.Logger.Debug().Msg("TelegramController: CreateNewUserForNowUserActive validate body")
	helpers.LogDebug("CreateNewUserForNowUserActive", "TelegramController: CreateNewUserForNowUserActive validate body", nil, c)
	var payload payloads.CreateNewTelegramPayload
	if err := helpers.ValidateBody(&payload, c); err != nil {
		// helpers.Logger.Error().Err(err).Msg("TelegramController: CreateNewUserForNowUserActive Error validating body")
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
		return helpers.ResponseErrorInternal(c, err)
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
		return helpers.ResponseErrorBadRequest(c, "Invalid telegram ID", nil)
	}
	tx := database.ClientPostgres.Begin()
	defer tx.Rollback()

	helpers.MyLogger("debug", "TelegramAccountLink", "DeleteByTelegramId", "controller", "start calling service for delete telegram user by telegram ID", nil, c)
	if err := t.TelegramService.DeleteByTelegramID(telegramIDInt64, c, tx); err != nil {
		return helpers.ResponseErrorInternal(c, err)
	}
	return nil
}

// find by user id
func (t *TelegramController) FindByUserID(c *fiber.Ctx) error {
	helpers.MyLogger("debug", "TelegramAccountLink", "FindByUserId", "controller", "start find telegram user by user ID", nil, c)
	userID := helpers.GetCurrentUserID(c)
	// parse userID to uint
	// userIDInt, err := strconv.ParseUint(userID, 10, 32)
	// if err != nil {
	// 	helpers.MyLogger("error", "TelegramAccountLink", "FindByUserId", "controller", "error parse user ID", map[string]interface{}{
	// 		"error": err.Error(),
	// 	}, c)
	// }
	// tx := database.ClientPostgres.Begin()
	// defer tx.Rollback()

	// helpers.MyLogger("debug", "TelegramAccountLink", "FindByUserId", "controller", "start calling service for find telegram user by user ID", map[string]interface{}{
	// 	"user_id":      userIDInt,
	// 	"user_id_uint": uint(userIDInt),
	// 	"user_id_str":  userID,
	// }, c)
	if err := t.TelegramService.FindByUserID(userID, c, database.ClientPostgres); err != nil {
		helpers.MyLogger("error", "TelegramAccountLink", "FindByUserId", "controller", "error find telegram user by user ID", map[string]interface{}{
			"error": err.Error(),
		}, c)
		return helpers.ResponseErrorInternal(c, err)
	}
	return nil
}