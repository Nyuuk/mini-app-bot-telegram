package repositories

import (
	"time"

	"github.com/Nyuuk/mini-app-bot-telegram/backend/app/entities"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type OvertimeRepository struct{}

// GetTelegramUserIDByTelegramID retrieves telegram_user.id by telegram_id
func (o *OvertimeRepository) GetTelegramUserIDByTelegramID(telegramID int64, c *fiber.Ctx, tx *gorm.DB) (uint, error) {
	var telegramUser entities.TelegramUser
	err := tx.WithContext(c.Context()).
		Where("telegram_id = ?", telegramID).
		First(&telegramUser).Error
	if err != nil {
		return 0, err
	}
	return telegramUser.ID, nil
}

// CreateNewRecordOvertime creates a new overtime record
func (o *OvertimeRepository) CreateNewRecordOvertime(payload *entities.Overtime, c *fiber.Ctx, tx *gorm.DB) error {
	err := tx.WithContext(c.Context()).Create(&payload).Error
	if err != nil {
		return err
	}
	return nil
}

// GetAllRecordOvertimeByTelegramID retrieves all overtime records by telegram ID
func (o *OvertimeRepository) GetAllRecordOvertimeByTelegramID(telegramID int64, overtimes *[]entities.Overtime, c *fiber.Ctx, tx *gorm.DB) error {
	err := tx.WithContext(c.Context()).
		Preload("User").
		Preload("TelegramUser").
		Joins("JOIN telegram_users ON telegram_users.id = overtimes.telegram_user_id").
		Where("telegram_users.telegram_id = ?", telegramID).
		Order("overtimes.date DESC, overtimes.created_at DESC").
		Find(&overtimes).Error
	if err != nil {
		return err
	}
	return nil
}

// GetRecordByDateByTelegramID retrieves overtime record by specific date and telegram ID
func (o *OvertimeRepository) GetRecordByDateByTelegramID(telegramID int64, date time.Time, overtime *[]entities.Overtime, c *fiber.Ctx, tx *gorm.DB) error {
	// Extract date string to avoid timezone conversion issues
	dateStr := date.Format("2006-01-02")

	err := tx.WithContext(c.Context()).
		Preload("User").
		Preload("TelegramUser").
		Joins("JOIN telegram_users ON telegram_users.id = overtimes.telegram_user_id").
		Where("telegram_users.telegram_id = ? AND DATE(overtimes.date) = ?", telegramID, dateStr).
		Find(&overtime).Error
	if err != nil {
		return err
	}
	return nil
}

// GetRecordBetweenDateByTelegramId retrieves overtime records between two dates for a telegram ID
func (o *OvertimeRepository) GetRecordBetweenDateByTelegramId(telegramID int64, startDate time.Time, endDate time.Time, overtimes *[]entities.Overtime, c *fiber.Ctx, tx *gorm.DB) error {
	// Extract date strings to avoid timezone conversion issues
	startDateStr := startDate.Format("2006-01-02")
	endDateStr := endDate.Format("2006-01-02")

	err := tx.WithContext(c.Context()).
		Preload("User").
		Preload("TelegramUser").
		Joins("JOIN telegram_users ON telegram_users.id = overtimes.telegram_user_id").
		Where("telegram_users.telegram_id = ? AND DATE(overtimes.date) BETWEEN ? AND ?", telegramID, startDateStr, endDateStr).
		Order("overtimes.date DESC, overtimes.created_at DESC").
		Find(&overtimes).Error
	if err != nil {
		return err
	}
	return nil
}

// GetRecordByID retrieves overtime record by ID
func (o *OvertimeRepository) GetRecordByID(id uint, overtime *entities.Overtime, c *fiber.Ctx, tx *gorm.DB) error {
	err := tx.WithContext(c.Context()).
		Preload("User").
		Preload("TelegramUser").
		Where("id = ?", id).
		First(&overtime).Error
	if err != nil {
		return err
	}
	return nil
}

// UpdateRecordOvertime updates an existing overtime record
func (o *OvertimeRepository) UpdateRecordOvertime(id uint, payload *entities.Overtime, c *fiber.Ctx, tx *gorm.DB) error {
	err := tx.WithContext(c.Context()).
		Model(&entities.Overtime{}).
		Where("id = ?", id).
		Updates(payload).Error
	if err != nil {
		return err
	}
	return nil
}

// DeleteRecordOvertime deletes an overtime record
func (o *OvertimeRepository) DeleteRecordOvertime(id uint, c *fiber.Ctx, tx *gorm.DB) error {
	err := tx.WithContext(c.Context()).
		Where("id = ?", id).
		Delete(&entities.Overtime{}).Error
	if err != nil {
		return err
	}
	return nil
}
