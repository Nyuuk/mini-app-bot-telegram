package controllers

import (
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

