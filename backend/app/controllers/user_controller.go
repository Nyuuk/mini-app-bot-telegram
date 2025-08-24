package controllers

import (
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
