package payloads

import (
	"time"

	"github.com/go-playground/validator/v10"
)

type CreateUserPayload struct {
	Username string `json:"username" validate:"required,min=3,max=20"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

func (p *CreateUserPayload) CustomErrorsMessage(validationErrors validator.ValidationErrors) []map[string]string {
	data := []map[string]string{}
	for _, err := range validationErrors {
		field := err.Field()
		switch field {
		case "Username":
			data = append(data, map[string]string{"username": "Username must be at least 3 characters"})
		case "Email":
			data = append(data, map[string]string{"email": "Email is required"})
		case "Password":
			data = append(data, map[string]string{"password": "Password must be at least 8 characters"})
		}
	}
	return data
}

type CreateApiKeyPayload struct {
	Description string     `json:"description" validate:"min=3,max=255"`
	IsActive    bool       `json:"is_active" validate:"required"`
	ExpiredAt   *time.Time `json:"expired_at" validate:"required"`
}

func (p *CreateApiKeyPayload) CustomErrorsMessage(errors validator.ValidationErrors) []map[string]string {
	var errorMessages []map[string]string
	for _, err := range errors {
		field := err.Field()
		switch field {
		case "Description":
			errorMessages = append(errorMessages, map[string]string{"description": "Description must be at least 3 characters"})
		case "IsActive":
			errorMessages = append(errorMessages, map[string]string{"is_active": "IsActive is required"})
		case "ExpiredAt":
			errorMessages = append(errorMessages, map[string]string{"expired_at": "ExpiredAt is required and must be in the future"})
		}
	}
	return errorMessages
}

type ResponseCreateApiKeyPayload struct {
	Description string     `json:"description"`
	IsActive    bool       `json:"is_active"`
	ExpiredAt   *time.Time `json:"expired_at"`
	APIKey      string     `json:"api_key"`
}
