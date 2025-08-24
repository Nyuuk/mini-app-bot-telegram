package payloads

import (
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
		data = append(data, map[string]string{
			"field":   err.Field(),
			"tag":     err.Tag(),
			"message": err.Error(),
		})
	}
	return data
}
