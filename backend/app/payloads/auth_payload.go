package payloads

import (
	"time"

	"github.com/Nyuuk/mini-app-bot-telegram/backend/app/entities"
	"github.com/go-playground/validator/v10"
)

type LoginPayload struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (p *LoginPayload) CustomErrorsMessage(errors validator.ValidationErrors) []map[string]string {
	var errorMessages []map[string]string
	for _, err := range errors {
		field := err.Field()
		switch field {
		case "Username":
			errorMessages = append(errorMessages, map[string]string{"username": "Username is required"})
		case "Password":
			errorMessages = append(errorMessages, map[string]string{"password": "Password is required"})
		}
	}
	return errorMessages
}

type RegisterPayload struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

func (p *RegisterPayload) CustomErrorsMessage(errors validator.ValidationErrors) []map[string]string {
	var errorMessages []map[string]string
	for _, err := range errors {
		field := err.Field()
		switch field {
		case "Username":
			errorMessages = append(errorMessages, map[string]string{"username": "Username must be between 3-50 characters"})
		case "Email":
			errorMessages = append(errorMessages, map[string]string{"email": "Valid email is required"})
		case "Password":
			errorMessages = append(errorMessages, map[string]string{"password": "Password must be at least 6 characters"})
		}
	}
	return errorMessages
}

type ResponseDetailMe struct {
	User     entities.User `json:"user"`
	AuthInfo AuthInfo      `json:"auth_info"`
}

type AuthInfo struct {
	UserID   uint      `json:"user_id"`
	AuthType string    `json:"auth_type"`
	ExpireAt time.Time `json:"expire_at"`
}
