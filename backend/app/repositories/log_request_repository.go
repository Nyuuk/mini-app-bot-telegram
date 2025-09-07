package repositories

import (
	"github.com/Nyuuk/mini-app-bot-telegram/backend/app/entities"
	"github.com/Nyuuk/mini-app-bot-telegram/backend/app/pkg/database"
	"github.com/Nyuuk/mini-app-bot-telegram/backend/app/pkg/helpers"
	"github.com/gofiber/fiber/v2"
)

type LogRequestRepository struct{}

func (r *LogRequestRepository) CreateLogRequest(log entities.LogRequest) (uint, error) {
	err := database.ClientPostgres.Create(&log).Error // &log is a pointer to the value
	return log.ID, err
}
func (r *LogRequestRepository) UpdateLogRequestEventByID(ID uint, event string) error {
	return database.ClientPostgres.Model(&entities.LogRequest{}).Where("id = ?", ID).Update("event", event).Error
}

func UpdateLogRequest(event string, c *fiber.Ctx) error {
	logRequestID := helpers.GetLogRequestId(c)
	repo := LogRequestRepository{}
	return repo.UpdateLogRequestEventByID(logRequestID, event)
}