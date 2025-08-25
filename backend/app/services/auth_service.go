package services

import (
	"strconv"
	"time"

	"github.com/Nyuuk/mini-app-bot-telegram/backend/app/entities"
	"github.com/Nyuuk/mini-app-bot-telegram/backend/app/payloads"
	"github.com/Nyuuk/mini-app-bot-telegram/backend/app/pkg/helpers"
	"github.com/Nyuuk/mini-app-bot-telegram/backend/app/repositories"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	ApiKeyRepository repositories.ApiKeyRepository
	UserRepository   repositories.UserRepository
}

func (s *AuthService) Login(c *fiber.Ctx, tx *gorm.DB, payload *payloads.LoginPayload) error {
	user := entities.User{}
	if err := s.UserRepository.FindByUsername(payload.Username, &user, tx); err != nil {
		return helpers.ResponseErrorBadRequest(c, "user or password is invalid", nil)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(payload.Password)); err != nil {
		return helpers.ResponseErrorBadRequest(c, "user or password is invalid", nil)
	}

	// generate jwt token
	expireAt := time.Now().Add(time.Hour * 24)
	userIdToString := strconv.Itoa(int(user.ID))
	log.Info("userIdToString: ", userIdToString)
	token, err := helpers.GenerateJWT(userIdToString, expireAt.Unix())
	if err != nil {
		return helpers.ResponseErrorInternal(c, err)
	}

	return helpers.Response(c, fiber.StatusOK, "Login successful", fiber.Map{
		"token":     token,
		"expire_at": expireAt.Format("2006-01-02 15:04:05"),
	})
}

func (s *AuthService) ValidateApiKey(apiKeyEntity *entities.APIKey, apiKey string, c *fiber.Ctx, tx *gorm.DB) bool {
	log.Info("Validating API key: ", apiKey)
	if err := s.ApiKeyRepository.FindByApiKey(apiKey, apiKeyEntity, tx); err != nil {
		log.Error("Error validating API key: ", err)
		return false
	}

	if err := bcrypt.CompareHashAndPassword([]byte(apiKeyEntity.APIKey), []byte(apiKey)); err != nil {
		log.Error("Error comparing hash and password: ", err)
		return false
	}

	log.Info("API key validated successfully ", apiKeyEntity.User)
	return true
}

func (s *AuthService) GetUserByApiKey(c *fiber.Ctx, tx *gorm.DB) error {
	log.Info("Getting user by API key")
	apiKey := c.Get("X-API-Key")
	var apiKeyEntity entities.APIKey
	if !s.ValidateApiKey(&apiKeyEntity, apiKey, c, tx) {
		return helpers.ResponseErrorBadRequest(c, "API key is invalid", nil)
	}
	return helpers.Response(c, fiber.StatusOK, "User retrieved successfully", apiKeyEntity.User)
}
