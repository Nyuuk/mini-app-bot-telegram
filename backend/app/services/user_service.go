package services

import (
	"github.com/Nyuuk/mini-app-bot-telegram/backend/app/entities"
	"github.com/Nyuuk/mini-app-bot-telegram/backend/app/payloads"
	"github.com/Nyuuk/mini-app-bot-telegram/backend/app/pkg/helpers"
	"github.com/Nyuuk/mini-app-bot-telegram/backend/app/repositories"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

type UserService struct {
	UserRepository repositories.UserRepository
}

func (u *UserService) GetUser(c *fiber.Ctx, tx *gorm.DB) error {
	id := c.Params("id")
	if id == "" {
		return helpers.ResponseErrorBadRequest(c, "ID is required", nil)
	}

	var user entities.User
	if err := u.UserRepository.FindByID(id, &user, tx); err != nil {
		return err
	}

	return helpers.Response(c, fiber.StatusOK, "User retrieved successfully", user)
}

func (u *UserService) CreateUser(payload *payloads.CreateUserPayload, c *fiber.Ctx, tx *gorm.DB) error {
	log.Debug("CreateUser service")
	var user entities.User
	user.Username = payload.Username
	user.Email = payload.Email
	user.PasswordHash = helpers.HashPassword(payload.Password)

	log.Debug("CreateUser service: ", user)

	if err := u.UserRepository.CreateUser(&user, tx, c); err != nil {
		log.Error("CreateUser service: ", err)
		if helpers.IsDuplicateKeyError(err) {
			return helpers.ResponseErrorBadRequest(c, "Username or email already exists", nil)
		}
		return err
	}

	tx.Commit()
	log.Debug("CreateUser service: ", user)
	return helpers.Response(c, fiber.StatusCreated, "User created successfully", user)

}
