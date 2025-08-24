package controllers

import (
	"github.com/Nyuuk/mini-app-bot-telegram/backend/app/pkg/database"
	"github.com/Nyuuk/mini-app-bot-telegram/backend/app/pkg/helpers"
	"github.com/Nyuuk/mini-app-bot-telegram/backend/app/services"
	"github.com/gofiber/fiber/v2"
)

type AuthController struct {
	AuthService services.AuthService
}

func (a *AuthController) CheckUserApiKey(c *fiber.Ctx) error {
	apiKey := c.Get("X-API-KEY")
	if apiKey == "" {
		return helpers.ResponseErrorBadRequest(c, "API key is required", nil)
	}

	tx := database.ClientPostgres
	if err := a.AuthService.GetUserByApiKey(c, tx); err != nil {
		return helpers.ResponseErrorInternal(c, err)
	}

	return nil

}
