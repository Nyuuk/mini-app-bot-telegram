package helpers

import (
	"crypto/rand"
	"encoding/hex"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) string {
	// log.Debug("HashPassword: ", password)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	// log.Debug("HashPassword: ", hashedPassword)
	if err != nil {
		LogError(err, "HashPassword", "HashPassword", nil, nil)
		// log.Error("HashPassword: ", err)
		return ""
	}
	return string(hashedPassword)
}

func VerifyPassword(password, hashedPassword string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) == nil
}

func GenerateAPIKey(length int) (string, error) {
	// length = jumlah byte random
	// hasil akhir jadi string hex dengan panjang 2x (karena tiap byte jadi 2 char)
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
