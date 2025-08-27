package middlewares

import (
	"strconv"
	"strings"
	"time"

	"github.com/Nyuuk/mini-app-bot-telegram/backend/app/entities"
	"github.com/Nyuuk/mini-app-bot-telegram/backend/app/pkg/database"
	"github.com/Nyuuk/mini-app-bot-telegram/backend/app/pkg/helpers"
	"github.com/Nyuuk/mini-app-bot-telegram/backend/app/repositories"
	"github.com/Nyuuk/mini-app-bot-telegram/backend/app/services"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

// AuthMiddleware mendukung autentikasi via API Key atau JWT Token
func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Cek API Key terlebih dahulu
		apiKey := c.Get("X-API-Key")
		if apiKey != "" {
			return validateApiKey(c, apiKey)
		}

		// Jika tidak ada API Key, cek JWT Token
		authHeader := c.Get("Authorization")
		if authHeader != "" {
			return validateJWT(c, authHeader)
		}

		// Jika keduanya tidak ada
		return helpers.Response(c, fiber.StatusUnauthorized, "Authentication required", nil)
	}
}

// validateApiKey memvalidasi API Key
func validateApiKey(c *fiber.Ctx, apiKey string) error {
	authService := services.AuthService{}
	tx := database.ClientPostgres
	apiKeyEntity := entities.APIKey{}

	if !authService.ValidateApiKey(&apiKeyEntity, apiKey, c, tx) {
		return helpers.Response(c, fiber.StatusUnauthorized, "Invalid API key", nil)
	}

	// Set user info ke context
	c.Locals("user_id", apiKeyEntity.UserID)
	c.Locals("user", apiKeyEntity.User)
	c.Locals("auth_type", "api_key")
	c.Locals("expired_at", apiKeyEntity.ExpiredAt.Unix())
	c.Locals("api_key_entity", apiKeyEntity)

	// log.Info("API Key authentication successful for user: ", apiKeyEntity.UserID)
	helpers.LogAuth("api_key_authentication_success", strconv.Itoa(int(apiKeyEntity.UserID)), true, map[string]interface{}{
		"user":         apiKeyEntity.User,
		"api_key_id":   apiKeyEntity.ID,
		"api_key_hash": apiKeyEntity.APIKey, // Log hash instead of plain API key
		"ip_address":   c.IP(),
		"user_agent":   c.Get("User-Agent"),
		"path":         c.Path(),
	})
	return c.Next()
}

// validateJWT memvalidasi JWT Token
func validateJWT(c *fiber.Ctx, authHeader string) error {
	token := strings.Replace(authHeader, "Bearer ", "", 1)
	if token == "" {
		return helpers.Response(c, fiber.StatusUnauthorized, "Invalid token format", nil)
	}

	bodyJWT, err := helpers.VerifyJWT(token, helpers.GetEnv("JWT_SECRET_KEY", "hello_world"))
	if err != nil {
		return helpers.Response(c, fiber.StatusForbidden, "Invalid token", nil)
	}

	// Validasi expire time
	if bodyJWT.ExpireAt < time.Now().Unix() {
		log.Info("Token expired: ", time.Unix(bodyJWT.ExpireAt, 0).Format("2006-01-02 15:04:05"))
		return helpers.Response(c, fiber.StatusForbidden, "Token expired", nil)
	}

	// Validasi user ID
	if bodyJWT.UserID == 0 {
		return helpers.Response(c, fiber.StatusForbidden, "Invalid token payload", nil)
	}

	// Validasi user di database
	userRepository := repositories.UserRepository{}
	tx := database.ClientPostgres
	user := entities.User{}

	// log.Info("Checking user JWT in database: ", bodyJWT.UserID)
	helpers.Logger.Info().Uint("user_id", bodyJWT.UserID).Msg("Checking user JWT in database")
	if err := userRepository.FindByID(bodyJWT.UserID, &user, tx); err != nil {
		return helpers.Response(c, fiber.StatusForbidden, "User not found", nil)
	}

	// Set user info ke context
	c.Locals("user_id", user.ID)
	c.Locals("user", user)
	c.Locals("auth_type", "jwt")
	c.Locals("expired_at", time.Unix(bodyJWT.ExpireAt, 0))
	c.Locals("expire_at", bodyJWT.ExpireAt)

	// log.Info("JWT authentication successful for user: ", bodyJWT.UserID)
	helpers.LogAuth("jwt_authentication_success", strconv.Itoa(int(bodyJWT.UserID)), true, map[string]interface{}{
		"user":          user,
		"jwt_user_id":   bodyJWT.UserID,
		"jwt_expire_at": time.Unix(bodyJWT.ExpireAt, 0).Format("2006-01-02 15:04:05"),
		"ip_address":    c.IP(),
		"user_agent":    c.Get("User-Agent"),
		"path":          c.Path(),
	})
	return c.Next()
}

// OptionalAuthMiddleware untuk endpoint yang bisa diakses tanpa auth
func OptionalAuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Cek API Key
		apiKey := c.Get("X-API-Key")
		if apiKey != "" {
			return validateApiKey(c, apiKey)
		}

		// Cek JWT Token
		authHeader := c.Get("Authorization")
		if authHeader != "" {
			return validateJWT(c, authHeader)
		}

		// Jika tidak ada auth, lanjutkan tanpa user info
		c.Locals("auth_type", "none")
		return c.Next()
	}
}
