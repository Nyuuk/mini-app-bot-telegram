package payloads

import (
	"github.com/go-playground/validator/v10"
)

type CreateNewRecordOvertime struct {
	TelegramID    int64   `json:"telegram_id" validate:"required"`
	Date          string  `json:"date" validate:"required"`       // Format: "2006-01-02" or "2006-01-02T15:04:05"
	TimeStart     string  `json:"time_start" validate:"required"` // Format: "2006-01-02T15:04:05"
	TimeStop      string  `json:"time_stop" validate:"required"`  // Format: "2006-01-02T15:04:05"
	BreakDuration float64 `json:"break_duration" validate:"gte=0"`
	Duration      float64 `json:"duration" validate:"required,gt=0"`
	Description   string  `json:"description" validate:"omitempty,min=3,max=255"`
	Category      string  `json:"category" validate:"omitempty,min=3,max=255"`
}

type GetRecordByDateRequest struct {
	TelegramID int64  `json:"telegram_id" validate:"required"`
	Date       string `json:"date" validate:"required" example:"2024-01-15"`
}

type GetRecordBetweenDateRequest struct {
	TelegramID int64  `json:"telegram_id" validate:"required"`
	StartDate  string `json:"start_date" validate:"required" example:"2024-01-01"`
	EndDate    string `json:"end_date" validate:"required" example:"2024-01-31"`
}

func (p *CreateNewRecordOvertime) CustomErrorsMessage(errors validator.ValidationErrors) []map[string]string {
	var errorMessages []map[string]string
	for _, err := range errors {
		field := err.Field()
		switch field {
		case "TelegramID":
			errorMessages = append(errorMessages, map[string]string{"telegram_id": "Telegram ID is required"})
		case "Date":
			errorMessages = append(errorMessages, map[string]string{"date": "Date is required"})
		case "TimeStart":
			errorMessages = append(errorMessages, map[string]string{"time_start": "Time start is required"})
		case "TimeStop":
			errorMessages = append(errorMessages, map[string]string{"time_stop": "Time stop is required"})
		case "BreakDuration":
			errorMessages = append(errorMessages, map[string]string{"break_duration": "Break duration must be greater than or equal to 0"})
		case "Duration":
			errorMessages = append(errorMessages, map[string]string{"duration": "Duration is required and must be greater than 0"})
		case "Description":
			errorMessages = append(errorMessages, map[string]string{"description": "Description must be at least 3 characters"})
		case "Category":
			errorMessages = append(errorMessages, map[string]string{"category": "Category must be at least 3 characters"})
		}
	}
	return errorMessages
}

func (p *GetRecordByDateRequest) CustomErrorsMessage(errors validator.ValidationErrors) []map[string]string {
	var errorMessages []map[string]string
	for _, err := range errors {
		field := err.Field()
		switch field {
		case "TelegramID":
			errorMessages = append(errorMessages, map[string]string{"telegram_id": "Telegram ID is required"})
		case "Date":
			errorMessages = append(errorMessages, map[string]string{"date": "Date is required"})
		}
	}
	return errorMessages
}

func (p *GetRecordBetweenDateRequest) CustomErrorsMessage(errors validator.ValidationErrors) []map[string]string {
	var errorMessages []map[string]string
	for _, err := range errors {
		field := err.Field()
		switch field {
		case "TelegramID":
			errorMessages = append(errorMessages, map[string]string{"telegram_id": "Telegram ID is required"})
		case "StartDate":
			errorMessages = append(errorMessages, map[string]string{"start_date": "Start date is required"})
		case "EndDate":
			errorMessages = append(errorMessages, map[string]string{"end_date": "End date is required"})
		}
	}
	return errorMessages
}
