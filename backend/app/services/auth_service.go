package services

import (
	"strconv"
	"time"

	"github.com/Nyuuk/mini-app-bot-telegram/backend/app/entities"
	"github.com/Nyuuk/mini-app-bot-telegram/backend/app/payloads"
	"github.com/Nyuuk/mini-app-bot-telegram/backend/app/pkg/helpers"
	"github.com/Nyuuk/mini-app-bot-telegram/backend/app/repositories"
	"github.com/gofiber/fiber/v2"
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
	expiredFromEnv := helpers.GetEnv("JWT_EXPIRATION", "1")
	expiredFromInt, err := strconv.Atoi(expiredFromEnv)
	if err != nil {
		return helpers.ResponseErrorInternal(c, err)
	}
	expireAt := time.Now().Add(time.Hour * time.Duration(expiredFromInt))
	userIdToString := strconv.Itoa(int(user.ID))
	helpers.Logger.Info().Str("user_id", userIdToString).Msg("Generating JWT token")
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
	helpers.Logger.Info().Str("api_key", apiKeyEntity.APIKey).Msg("Validating API key")
	if err := s.ApiKeyRepository.FindByApiKey(apiKey, apiKeyEntity, tx); err != nil {
		helpers.Logger.Error().Err(err).Str("api_key", apiKeyEntity.APIKey).Msg("Error validating API key")
		return false
	}

	// if err := bcrypt.CompareHashAndPassword([]byte(apiKeyEntity.APIKey), []byte(apiKey)); err != nil {
	// 	helpers.Logger.Error().Err(err).Str("api_key", apiKeyEntity.APIKey).Msg("Error comparing hash and password")
	// 	return false
	// }

	helpers.Logger.Info().Str("api_key", apiKeyEntity.APIKey).Str("user_id", strconv.Itoa(int(apiKeyEntity.UserID))).Msg("API key validated successfully")
	return true
}

func (s *AuthService) GetUserByApiKey(c *fiber.Ctx, tx *gorm.DB) error {
	helpers.Logger.Info().Msg("Getting user by API key")
	apiKey := c.Get("X-API-Key")
	var apiKeyEntity entities.APIKey
	if !s.ValidateApiKey(&apiKeyEntity, apiKey, c, tx) {
		return helpers.ResponseErrorBadRequest(c, "API key is invalid", nil)
	}
	return helpers.Response(c, fiber.StatusOK, "User retrieved successfully", apiKeyEntity.User)
}
