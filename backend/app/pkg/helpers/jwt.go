package helpers

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

type JWTBody struct {
	UserID   string `json:"user_id"`
	ExpireAt int64  `json:"expire_at"`
}

func GenerateJWT(userID string, expireAt int64) (string, error) {
	secretKey := GetEnv("JWT_SECRET_KEY", "hello_world")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":   userID,
		"expire_at": expireAt,
	})

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyJWT(tokenString string, secretKey string) (JWTBody, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		return JWTBody{}, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return JWTBody{}, errors.New("invalid token claims")
	}

	return JWTBody{
		UserID:   claims["user_id"].(string),
		ExpireAt: int64(claims["expire_at"].(float64)),
	}, nil
}
