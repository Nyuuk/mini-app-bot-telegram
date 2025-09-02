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

// Login godoc
// @Summary User Login
// @Description Login with username and password to get JWT token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param loginPayload body payloads.LoginPayload true "Login credentials"
// @Success 200 {object} map[string]interface{} "Login successful"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /v1/auth/login [post]
func (a *AuthController) Login(c *fiber.Ctx) error {
	// Log login attempt
	helpers.LogAuth("login_attempt", "anonymous", false, map[string]interface{}{
		"ip_address": c.IP(),
		"user_agent": c.Get("User-Agent"),
		"path":       c.Path(),
	})

	payload := payloads.LoginPayload{}
	if err := helpers.ValidateBody(&payload, c); err != nil {
		// Log validation error
		// helpers.LogError(err, "login_validation_error", "anonymous", map[string]interface{}{
		// 	"ip_address": c.IP(),
		// 	"user_agent": c.Get("User-Agent"),
		// 	"username":   payload.Username,
		// })
		helpers.LogError(err, "login_validation_error", "AuthController: error when validating body", nil, c)
		return helpers.ResponseErrorBadRequest(c, "Invalid request body", err)
	}

	tx := database.ClientPostgres
	err := a.AuthService.Login(c, tx, &payload)
	if err != nil {
		// Log login failure
		helpers.LogAuth("login_failed", "anonymous", false, map[string]interface{}{
			"ip_address": c.IP(),
			"user_agent": c.Get("User-Agent"),
			"username":   payload.Username,
			"error":      err.Error(),
		})
		return helpers.ResponseErrorInternal(c, err)
	}

	// Log successful login
	helpers.LogAuth("login_success", "anonymous", true, map[string]interface{}{
		"ip_address": c.IP(),
		"user_agent": c.Get("User-Agent"),
		"username":   payload.Username,
	})

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
