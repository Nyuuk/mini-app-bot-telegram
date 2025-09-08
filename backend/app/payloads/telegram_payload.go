package payloads

import (
	"github.com/Nyuuk/mini-app-bot-telegram/backend/app/pkg/errors"
	"github.com/go-playground/validator/v10"
)

type CreateNewTelegramPayload struct {
	TelegramID int64  `json:"telegram_id" validate:"required"`
	Username   string `json:"username" validate:"required,min=3,max=50"`
	FirstName  string `json:"first_name" validate:"min=3,max=50"`
	LastName   string `json:"last_name" validate:"min=3,max=50"`
}

func (p *CreateNewTelegramPayload) CustomErrorsMessage(validationErrors validator.ValidationErrors) []map[string]string {
	var errorMessages []map[string]string
	for _, err := range validationErrors {
		field := err.Field()
		tag := err.Tag()
		
		switch field {
		case "TelegramID":
			if tag == "required" {
				errorMessages = append(errorMessages, map[string]string{
					"telegram_id": "Telegram ID is required. Please provide a valid Telegram user ID.",
				})
			}
		case "Username":
			if tag == "required" {
				errorMessages = append(errorMessages, map[string]string{
					"username": "Telegram username is required. Please provide your Telegram username.",
				})
			} else if tag == "min" || tag == "max" {
				errorMessages = append(errorMessages, map[string]string{
					"username": "Telegram username must be between 3 and 50 characters long.",
				})
			}
		case "FirstName":
			if tag == "min" || tag == "max" {
				errorMessages = append(errorMessages, map[string]string{
					"first_name": "First name must be between 3 and 50 characters long.",
				})
			}
		case "LastName":
			if tag == "min" || tag == "max" {
				errorMessages = append(errorMessages, map[string]string{
					"last_name": "Last name must be between 3 and 50 characters long.",
				})
			}
		}
	}
	
	// If no specific errors were matched, provide a generic validation error
	if len(errorMessages) == 0 {
		errorMessages = append(errorMessages, map[string]string{
			"validation_error": "Invalid payload. Please check your request data and try again.",
		})
	}
	
	return errorMessages
}

type UpdateTelegramPayload struct {
	Username   string `json:"username" validate:"required,min=3,max=50"`
	FirstName  string `json:"first_name" validate:"min=3,max=50"`
	LastName   string `json:"last_name" validate:"min=3,max=50"`
}

func (p *UpdateTelegramPayload) CustomErrorsMessage(validationErrors validator.ValidationErrors) []map[string]string {
	var errorMessages []map[string]string
	for _, err := range validationErrors {
		field := err.Field()
		tag := err.Tag()
		
		switch field {
		case "Username":
			if tag == "required" {
				errorMessages = append(errorMessages, map[string]string{
					"username": "Telegram username is required for update.",
				})
			} else if tag == "min" || tag == "max" {
				errorMessages = append(errorMessages, map[string]string{
					"username": "Telegram username must be between 3 and 50 characters long.",
				})
			}
		case "FirstName":
			if tag == "min" || tag == "max" {
				errorMessages = append(errorMessages, map[string]string{
					"first_name": "First name must be between 3 and 50 characters long.",
				})
			}
		case "LastName":
			if tag == "min" || tag == "max" {
				errorMessages = append(errorMessages, map[string]string{
					"last_name": "Last name must be between 3 and 50 characters long.",
				})
			}
		}
	}
	
	// If no specific errors were matched, provide a generic validation error
	if len(errorMessages) == 0 {
		errorMessages = append(errorMessages, map[string]string{
			"validation_error": "Invalid update payload. Please check your request data and try again.",
		})
	}
	
	return errorMessages
}

// TelegramErrorResponse creates a standardized error response for Telegram operations
func TelegramErrorResponse(errorCode, message string, details any) errors.TelegramError {
	return errors.NewTelegramError(errorCode, message, details)
}