package controllers

import (
	"github.com/Nyuuk/mini-app-bot-telegram/backend/app/payloads"
	"github.com/Nyuuk/mini-app-bot-telegram/backend/app/pkg/database"
	"github.com/Nyuuk/mini-app-bot-telegram/backend/app/pkg/helpers"
	"github.com/Nyuuk/mini-app-bot-telegram/backend/app/services"
	"github.com/gofiber/fiber/v2"
)

type UserController struct {
	UserService services.UserService
}

// GetUserById godoc
// @Summary Get User by ID
// @Description Get user information by user ID
// @Tags Users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} map[string]interface{} "User retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid user ID"
// @Failure 404 {object} map[string]interface{} "User not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security ApiKeyAuth
// @Security BearerAuth
// @Router /v1/user/{id} [get]
func (u *UserController) GetUserById(c *fiber.Ctx) error {
	tx := database.ClientPostgres
	if err := u.UserService.GetUserById(c, tx); err != nil {
		return helpers.ResponseErrorInternal(c, err)
	}

	return nil
}

func (u *UserController) GetAllUsers(c *fiber.Ctx) error {
	tx := database.ClientPostgres
	if err := u.UserService.GetAllUsers(c, tx); err != nil {
		return helpers.ResponseErrorInternal(c, err)
	}

	return nil
}

func (u *UserController) GetDetailMe(c *fiber.Ctx) error {
	user := helpers.GetCurrentUser(c)
	var authInfo payloads.AuthInfo
	var responseDetailMe payloads.ResponseDetailMe
	authInfo.AuthType = helpers.GetAuthType(c)
	authInfo.ExpireAt = helpers.GetExpireAt(c)
	authInfo.UserID = helpers.GetCurrentUserID(c)
	responseDetailMe.User = user
	responseDetailMe.AuthInfo = authInfo
	return helpers.Response(c, fiber.StatusOK, "User retrieved successfully", responseDetailMe)
}

func (u *UserController) GetApiKeyFromUserActive(c *fiber.Ctx) error {
	helpers.LogDebug("GetApiKeyFromUserActive", "UserController: calling service GetApiKeyFromUserActive", nil, c)
	tx := database.ClientPostgres
	if err := u.UserService.GetApiKeyFromUserActive(c, tx); err != nil {
		// log.Error("GetApiKeyFromUserActive Controller: error calling service GetApiKeyFromUserActive: ", err)
		helpers.LogError(err, "GetApiKeyFromUserActive", "UserController: error calling service GetApiKeyFromUserActive", nil, c)
		return helpers.ResponseErrorInternal(c, err)
	}
	return nil
}

func (u *UserController) DeleteUserById(c *fiber.Ctx) error {
	tx := database.ClientPostgres
	if err := u.UserService.DeleteUserById(c, tx); err != nil {
		return helpers.ResponseErrorInternal(c, err)
	}
	return nil
}

// CreateUser godoc
// @Summary Create New User (Register)
// @Description Create a new user account
// @Tags Users
// @Accept json
// @Produce json
// @Param createUserPayload body payloads.CreateUserPayload true "User registration data"
// @Success 201 {object} map[string]interface{} "User created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /v1/auth/register [post]
func (u *UserController) CreateUser(c *fiber.Ctx) error {
	tx := database.ClientPostgres.Begin()
	defer tx.Rollback()

	payload := payloads.CreateUserPayload{}
	if err := helpers.ValidateBody(&payload, c); err != nil {
		// log.Error("ValidateBody: ", err)
		helpers.LogError(err, "CreateUser", "UserController: error when validating body", nil, c)
		if customErr, ok := err.(helpers.Error); ok {
			return helpers.ResponseErrorBadRequest(c, customErr.Message, customErr.Data)
		}
	}

	// log.Debug("Calling service.CreateUser with payload: ", payload)
	helpers.LogDebug("CreateUser", "UserController: calling service.CreateUser with payload", nil, c)
	if err := u.UserService.CreateUser(&payload, c, tx); err != nil {
		tx.Rollback()
		// log.Error("CreateUser: ", err)
		helpers.LogError(err, "CreateUser", "UserController: error when calling service.CreateUser", nil, c)
		return helpers.ResponseErrorInternal(c, err)
	}

	return nil

}

func (u *UserController) CreateApiKey(c *fiber.Ctx) error {
	payload := payloads.CreateApiKeyPayload{}
	if err := helpers.ValidateBody(&payload, c); err != nil {
		// log.Error("ValidateBody Controller CreateApiKey: ", err)
		helpers.LogError(err, "CreateApiKey", "UserController: error when validating body", nil, c)
		return helpers.ResponseErrorBadRequest(c, "Invalid request body", err)
	}

	tx := database.ClientPostgres.Begin()
	defer tx.Rollback()

	// log.Debug("Calling service.CreateApiKey with payload: ", payload)

	// helpers.LogDebug("CreateApiKey", "UserController: calling service.CreateApiKey with payload", map[string]interface{}{
	// 	"payload.Description": payload.Description,
	// 	"payload.ExpiresAt":   payload.ExpiredAt,
	// 	"payload.IsActive":    payload.IsActive,
	// }, c)
	helpers.Logger.Debug().
		Interface("payload", payload).
		Str("type", "debug").
		Str("event", "CreateApiKey").
		Str("payload.Description", payload.Description).
		Str("payload.ExpiredAt", payload.ExpiredAt.Format("2006-01-02 15:04:05")).
		Bool("payload.IsActive", payload.IsActive).
		Str("function", "UserController").
		Msg("calling service.CreateApiKey with payload")
	if err := u.UserService.CreateApiKey(&payload, c, tx); err != nil {
		tx.Rollback()
		// log.Error("CreateApiKey Controller: ", err)
		helpers.LogError(err, "CreateApiKey", "UserController: error when calling service.CreateApiKey", nil, c)
		return helpers.ResponseErrorInternal(c, err)
	}

	return nil
}
