package payloads

import "github.com/go-playground/validator/v10"

type CreateNewTelegramPayload struct {
	TelegramID int64  `json:"telegram_id" validate:"required"`
	Username   string `json:"username" validate:"required,min=3,max=50"`
	FirstName  string `json:"first_name" validate:"min=3,max=50"`
	LastName   string `json:"last_name" validate:"min=3,max=50"`
}

func (p *CreateNewTelegramPayload) CustomErrorsMessage(errors validator.ValidationErrors) []map[string]string {
	var errorMessages []map[string]string
	for _, err := range errors {
		field := err.Field()
		switch field {
		case "TelegramID":
			errorMessages = append(errorMessages, map[string]string{"telegram_id": "Telegram ID is required"})
		case "Username":
			errorMessages = append(errorMessages, map[string]string{"username": "Username is required"})
		case "FirstName":
			errorMessages = append(errorMessages, map[string]string{"first_name": "First name must be at least 3 characters"})
		case "LastName":
			errorMessages = append(errorMessages, map[string]string{"last_name": "Last name must be at least 3 characters"})
		}
	}
	return errorMessages
}
