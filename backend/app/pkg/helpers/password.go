package helpers

import (
	"github.com/gofiber/fiber/v2/log"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) string {
	log.Debug("HashPassword: ", password)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	log.Debug("HashPassword: ", hashedPassword)
	if err != nil {
		log.Error("HashPassword: ", err)
		return ""
	}
	return string(hashedPassword)
}

func VerifyPassword(password, hashedPassword string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) == nil
}
