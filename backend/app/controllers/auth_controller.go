package controllers

import (
	"github.com/Nyuuk/mini-app-bot-telegram/backend/app/payloads"
	"github.com/Nyuuk/mini-app-bot-telegram/backend/app/pkg/database"
	"github.com/Nyuuk/mini-app-bot-telegram/backend/app/pkg/helpers"
	"github.com/Nyuuk/mini-app-bot-telegram/backend/app/services"
	"github.com/gofiber/fiber/v2"
)

type AuthController struct {
	AuthService services.AuthService
}

// Login untuk mendapatkan JWT token
func (a *AuthController) Login(c *fiber.Ctx) error {
	payload := payloads.LoginPayload{}
	if err := helpers.ValidateBody(&payload, c); err != nil {
		return helpers.ResponseErrorBadRequest(c, "Invalid request body", err)
	}

	tx := database.ClientPostgres
	err := a.AuthService.Login(c, tx, &payload)
	if err != nil {
		return helpers.ResponseErrorInternal(c, err)
	}

	return nil
}


// // GetUserApiKeys mendapatkan semua API key user yang login
// func (a *AuthController) GetUserApiKeys(c *fiber.Ctx) error {
// 	userID := helpers.GetCurrentUserID(c)
// 	tx := database.ClientPostgres

// 	apiKeys, err := a.AuthService.GetUserApiKeys(c, tx, userID)
// 	if err != nil {
// 		return helpers.ResponseErrorInternal(c, err)
// 	}

// 	return helpers.Response(c, fiber.StatusOK, "API keys retrieved successfully", apiKeys)
// }

// DeleteApiKey menghapus API key user
// func (a *AuthController) DeleteApiKey(c *fiber.Ctx) error {
// 	apiKeyID := c.Params("id")
// 	userID := helpers.GetCurrentUserID(c)

// 	tx := database.ClientPostgres.Begin()
// 	defer tx.Rollback()

// 	if err := a.AuthService.DeleteApiKey(c, tx, userID, apiKeyID); err != nil {
// 		return helpers.ResponseErrorInternal(c, err)
// 	}

// 	tx.Commit()
// 	return helpers.Response(c, fiber.StatusOK, "API key deleted successfully", nil)
// }
