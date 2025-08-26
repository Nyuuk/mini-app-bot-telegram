package services

import (
	"fmt"
	"strconv"

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
	// parse id to uint
	userID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return helpers.ResponseErrorBadRequest(c, "Invalid ID", nil)
	}
	if err := u.UserRepository.FindByID(uint(userID), &user, tx); err != nil {
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

func (u *UserService) CreateApiKey(payload *payloads.CreateApiKeyPayload, c *fiber.Ctx, tx *gorm.DB) error {
	log.Debug("CreateApiKey service")
	userID := helpers.GetCurrentUserID(c)

	log.Debug("CreateApiKey service: Finding user by ID: ", userID)
	var user entities.User
	if err := u.UserRepository.FindByID(userID, &user, tx); err != nil {
		return err
	}
	log.Debug("CreateApiKey service: User found: ", user)

	var apiKey entities.APIKey
	apiKey.UserID = userID
	apiKey.Description = payload.Description
	apiKey.IsActive = payload.IsActive
	apiKey.ExpiredAt = payload.ExpiredAt
	apiKeyGenerated, err := helpers.GenerateAPIKey(32)
	log.Debug("CreateApiKey service: Generated API key: ", apiKeyGenerated)
	if err != nil {
		return err
	}
	// encrypt api key
	// apiKey.APIKey = helpers.HashPassword(apiKeyGenerated)
	apiKey.APIKey = apiKeyGenerated

	log.Debug(fmt.Sprintf("CreateApiKey service: Creating API key: %s Description: %s ExpiredAt: %s IsActive: %t", apiKeyGenerated, payload.Description, payload.ExpiredAt, payload.IsActive))
	if err := u.UserRepository.CreateApiKey(&apiKey, tx, c); err != nil {
		log.Error("CreateApiKey service: Error creating API key: ", err)
		return err
	}

	log.Debug("CreateApiKey service: Committing transaction")
	// tx.Rollback()
	if err := tx.Commit().Error; err != nil {
		log.Error("CreateApiKey service: Error committing transaction: ", err)
		return err
	}

	var response payloads.ResponseCreateApiKeyPayload
	response.APIKey = apiKeyGenerated
	response.Description = apiKey.Description
	response.ExpiredAt = apiKey.ExpiredAt
	response.IsActive = apiKey.IsActive

	return helpers.Response(c, fiber.StatusOK, "API key created successfully", response)
}
