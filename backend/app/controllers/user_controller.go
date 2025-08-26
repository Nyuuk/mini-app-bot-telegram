package controllers

import (
	"github.com/Nyuuk/mini-app-bot-telegram/backend/app/entities"
	"github.com/Nyuuk/mini-app-bot-telegram/backend/app/payloads"
	"github.com/Nyuuk/mini-app-bot-telegram/backend/app/pkg/database"
	"github.com/Nyuuk/mini-app-bot-telegram/backend/app/pkg/helpers"
	"github.com/Nyuuk/mini-app-bot-telegram/backend/app/services"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type UserController struct {
	UserService services.UserService
}

func (u *UserController) GetUserById(c *fiber.Ctx) error {
	tx := database.ClientPostgres
	if err := u.UserService.GetUser(c, tx); err != nil {
		return helpers.ResponseErrorInternal(c, err)
	}

	return nil
}

func (u *UserController) CreateUser(c *fiber.Ctx) error {
	log.Info("CreateUser")
	tx := database.ClientPostgres.Begin()
	defer tx.Rollback()

	payload := payloads.CreateUserPayload{}
	if err := helpers.ValidateBody(&payload, c); err != nil {
		log.Error("ValidateBody: ", err)
		if customErr, ok := err.(helpers.Error); ok {
			return helpers.ResponseErrorBadRequest(c, customErr.Message, customErr.Data)
		}
	}

	log.Debug("Calling service.CreateUser with payload: ", payload)
	if err := u.UserService.CreateUser(&payload, c, tx); err != nil {
		tx.Rollback()
		log.Error("CreateUser: ", err)
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
	responseDetailMe.User = user.(entities.User)
	responseDetailMe.AuthInfo = authInfo
	return helpers.Response(c, fiber.StatusOK, "User retrieved successfully", responseDetailMe)
}

func (u *UserController) CreateApiKey(c *fiber.Ctx) error {
	payload := payloads.CreateApiKeyPayload{}
	if err := helpers.ValidateBody(&payload, c); err != nil {
		log.Error("ValidateBody Controller CreateApiKey: ", err)
		return helpers.ResponseErrorBadRequest(c, "Invalid request body", err)
	}

	tx := database.ClientPostgres.Begin()
	defer tx.Rollback()

	log.Debug("Calling service.CreateApiKey with payload: ", payload)
	if err := u.UserService.CreateApiKey(&payload, c, tx); err != nil {
		tx.Rollback()
		log.Error("CreateApiKey Controller: ", err)
		return helpers.ResponseErrorInternal(c, err)
	}

	return nil
}
