package helpers

import (
	"github.com/gofiber/fiber/v2"
)

// GetCurrentUserID mendapatkan user ID dari context
func GetCurrentUserID(c *fiber.Ctx) uint {
	if userID, ok := c.Locals("user_id").(uint); ok {
		return userID
	}
	return 0
}

// GetCurrentUser mendapatkan user object dari context
func GetCurrentUser(c *fiber.Ctx) interface{} {
	return c.Locals("user")
}

// GetAuthType mendapatkan tipe autentikasi yang digunakan
func GetAuthType(c *fiber.Ctx) string {
	if authType, ok := c.Locals("auth_type").(string); ok {
		return authType
	}
	return "none"
}

// IsAuthenticated mengecek apakah user sudah terautentikasi
func IsAuthenticated(c *fiber.Ctx) bool {
	return GetCurrentUserID(c) != 0
}

// RequireAuth middleware untuk memastikan user terautentikasi
func RequireAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if !IsAuthenticated(c) {
			return Response(c, fiber.StatusUnauthorized, "Authentication required", nil)
		}
		return c.Next()
	}
}
