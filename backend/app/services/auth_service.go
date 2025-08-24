package services

import (
	"github.com/Nyuuk/mini-app-bot-telegram/backend/app/entities"
	"github.com/Nyuuk/mini-app-bot-telegram/backend/app/pkg/helpers"
	"github.com/Nyuuk/mini-app-bot-telegram/backend/app/repositories"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	ApiKeyRepository repositories.ApiKeyRepository
}

func (s *AuthService) ValidateApiKey(apiKeyEntity *entities.APIKey, c *fiber.Ctx, tx *gorm.DB) bool {
	apiKey := c.Get("X-API-KEY")
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
	var apiKeyEntity entities.APIKey
	if !s.ValidateApiKey(&apiKeyEntity, c, tx) {
		return helpers.ResponseErrorBadRequest(c, "API key is invalid", nil)
	}
	return helpers.Response(c, fiber.StatusOK, "User retrieved successfully", apiKeyEntity.User)
}
