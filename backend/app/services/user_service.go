package services

import (
	"errors"
	"strconv"

	"github.com/Nyuuk/mini-app-bot-telegram/backend/app/entities"
	"github.com/Nyuuk/mini-app-bot-telegram/backend/app/payloads"
	"github.com/Nyuuk/mini-app-bot-telegram/backend/app/pkg/helpers"
	"github.com/Nyuuk/mini-app-bot-telegram/backend/app/repositories"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type UserService struct {
	UserRepository repositories.UserRepository
}

func (u *UserService) GetUserById(c *fiber.Ctx, tx *gorm.DB) error {
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
		if helpers.IsNotFoundError(err) {
			return helpers.ResponseErrorNotFound(c, nil)
		}
		return err
	}

	return helpers.Response(c, fiber.StatusOK, "User retrieved successfully", user)
}

func (u *UserService) GetAllUsers(c *fiber.Ctx, tx *gorm.DB) error {
	var users []entities.User
	if err := u.UserRepository.FindAll(&users, tx); err != nil {
		return err
	}
	return helpers.Response(c, fiber.StatusOK, "Users retrieved successfully", users)
}

func (u *UserService) GetApiKeyFromUserActive(c *fiber.Ctx, tx *gorm.DB) error {
	userID := helpers.GetCurrentUserID(c)
	helpers.LogDebug("GetApiKeyFromUserActive", "UserService: GetApiKeyFromUserActive", map[string]interface{}{
		"userID": userID,
	}, c)
	var user entities.User
	if err := u.UserRepository.FindByID(userID, &user, tx); err != nil {
		if helpers.IsNotFoundError(err) {
			return helpers.ResponseErrorNotFound(c, nil)
		}
		helpers.LogError(err, "GetApiKeyFromUserActive", "UserService: error finding user by ID", nil, c)
		return err
	}
	return helpers.Response(c, fiber.StatusOK, "API key retrieved successfully", user.APIKeys)
}

func (u *UserService) DeleteUserById(c *fiber.Ctx, tx *gorm.DB) error {
	id := c.Params("id")
	if id == "" {
		helpers.LogError(errors.New("ID is required"), "DeleteUserById", "UserService: ID is required", nil, c)
		return helpers.ResponseErrorBadRequest(c, "ID is required", nil)
	}

	var user entities.User
	userID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		helpers.LogError(err, "DeleteUserById", "UserService: error when parsing ID", nil, c)
		return helpers.ResponseErrorBadRequest(c, "Invalid ID", nil)
	}
	if err := u.UserRepository.FindByID(uint(userID), &user, tx); err != nil {
		if helpers.IsNotFoundError(err) {
			helpers.LogInfo("DeleteUserById", "UserService: user not found", nil, c)
			return helpers.ResponseErrorNotFound(c, nil)
		}
		helpers.LogError(err, "DeleteUserById", "UserService: error finding user by ID", nil, c)
		return err
	}

	helpers.LogDebug("DeleteUserById", "UserService: calling repository DeleteUserById", nil, c)
	if err := u.UserRepository.DeleteUserById(uint(userID), tx); err != nil {
		helpers.LogError(err, "DeleteUserById", "UserService: error deleting user by ID", nil, c)
		return err
	}

	helpers.LogInfo("DeleteUserById", "UserService: user deleted successfully", nil, c)
	return helpers.Response(c, fiber.StatusOK, "User deleted successfully", user)
}

func (u *UserService) CreateUser(payload *payloads.CreateUserPayload, c *fiber.Ctx, tx *gorm.DB) error {
	helpers.LogDebug("CreateUser", "UserService: CreateUser", map[string]interface{}{
		"payload": payload,
	}, c)
	var user entities.User
	user.Username = payload.Username
	user.Email = payload.Email
	user.PasswordHash = helpers.HashPassword(payload.Password)

	helpers.LogDebug("CreateUser", "UserService: user created successfully", map[string]interface{}{
		"user": user,
	}, c)

	if err := u.UserRepository.CreateUser(&user, tx, c); err != nil {
		helpers.LogError(err, "CreateUser", "UserService: error creating user", nil, c)
		if helpers.IsDuplicateKeyError(err) {
			helpers.LogInfo("CreateUser", "UserService: username or email already exists", nil, c)
			return helpers.ResponseErrorBadRequest(c, "Username or email already exists", nil)
		}
		helpers.LogError(err, "CreateUser", "UserService: error creating user", nil, c)
		return err
	}

	tx.Commit()
	helpers.LogDebug("CreateUser", "UserService: user created successfully", map[string]interface{}{
		"user": user,
	}, c)
	return helpers.Response(c, fiber.StatusCreated, "User created successfully", user)

}

func (u *UserService) CreateApiKey(payload *payloads.CreateApiKeyPayload, c *fiber.Ctx, tx *gorm.DB) error {
	helpers.LogDebug("CreateApiKey", "UserService: CreateApiKey", map[string]interface{}{
		"payload": payload,
	}, c)
	userID := helpers.GetCurrentUserID(c)

	helpers.LogDebug("CreateApiKey", "UserService: Finding user by ID", map[string]interface{}{
		"userID": userID,
	}, c)
	var user entities.User
	if err := u.UserRepository.FindByID(userID, &user, tx); err != nil {
		helpers.LogError(err, "CreateApiKey", "UserService: error finding user by ID", nil, c)
		return err
	}
	helpers.LogDebug("CreateApiKey", "UserService: User found", map[string]interface{}{
		"user": user,
	}, c)

	var apiKey entities.APIKey
	apiKey.UserID = userID
	apiKey.Description = payload.Description
	apiKey.IsActive = payload.IsActive
	apiKey.ExpiredAt = payload.ExpiredAt
	apiKeyGenerated, err := helpers.GenerateAPIKey(32)
	helpers.LogDebug("CreateApiKey", "UserService: Generated API key", map[string]interface{}{
		"apiKeyGenerated": apiKeyGenerated,
	}, c)
	if err != nil {
		helpers.LogError(err, "CreateApiKey", "UserService: error generating API key", nil, c)
		return err
	}
	// encrypt api key
	// apiKey.APIKey = helpers.HashPassword(apiKeyGenerated)
	apiKey.APIKey = apiKeyGenerated

	// log.Debug(fmt.Sprintf("CreateApiKey service: Creating API key: %s Description: %s ExpiredAt: %s IsActive: %t", apiKeyGenerated, payload.Description, payload.ExpiredAt, payload.IsActive))
	if err := u.UserRepository.CreateApiKey(&apiKey, tx, c); err != nil {
		helpers.LogError(err, "CreateApiKey", "UserService: error creating API key", nil, c)
		return err
	}

	helpers.LogDebug("CreateApiKey", "UserService: Committing transaction", nil, c)
	// tx.Rollback()
	if err := tx.Commit().Error; err != nil {
		helpers.LogError(err, "CreateApiKey", "UserService: error committing transaction", nil, c)
		return err
	}

	var response payloads.ResponseCreateApiKeyPayload
	response.APIKey = apiKeyGenerated
	response.Description = apiKey.Description
	response.ExpiredAt = apiKey.ExpiredAt
	response.IsActive = apiKey.IsActive

	return helpers.Response(c, fiber.StatusOK, "API key created successfully", response)
}
