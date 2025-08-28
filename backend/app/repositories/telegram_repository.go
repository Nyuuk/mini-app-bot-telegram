package repositories

import (
	"github.com/Nyuuk/mini-app-bot-telegram/backend/app/entities"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type TelegramRepository struct{}

func (t *TelegramRepository) Create(payload *entities.TelegramUser, c *fiber.Ctx, tx *gorm.DB) error {
	err := tx.WithContext(c.Context()).Create(&payload).Error
	if err != nil {
		return err
	}
	return nil
}

func (t *TelegramRepository) FindByUserID(userID uint, telegramUser *[]entities.TelegramUser, c *fiber.Ctx, tx *gorm.DB) error {
	// get all telegram user by user ID
	err := tx.WithContext(c.Context()).Where("user_id = ?", userID).Find(&telegramUser).Error
	if err != nil {
		return err
	}
	return nil
}

func (t *TelegramRepository) FindByID(ID uint, telegramUser *entities.TelegramUser, c *fiber.Ctx, tx *gorm.DB) error {
	err := tx.WithContext(c.Context()).Where("id = ?", ID).First(&telegramUser).Error
	if err != nil {
		return err
	}
	return nil
}

func (t *TelegramRepository) FindByTelegramID(telegramID int64, telegramUser *entities.TelegramUser, c *fiber.Ctx, tx *gorm.DB) error {
	err := tx.WithContext(c.Context()).Where("telegram_id = ?", telegramID).First(&telegramUser).Error
	if err != nil {
		return err
	}
	return nil
}

func (t *TelegramRepository) UpdateByTelegramID(telegramID int64, payload *entities.TelegramUser, c *fiber.Ctx, tx *gorm.DB) error {
	err := tx.WithContext(c.Context()).Model(&entities.TelegramUser{}).Where("telegram_id = ?", telegramID).Updates(payload).Error
	if err != nil {
		return err
	}
	return nil
}

func (t *TelegramRepository) UpdateByID(ID uint, payload *entities.TelegramUser, c *fiber.Ctx, tx *gorm.DB) error {
	err := tx.WithContext(c.Context()).Model(&entities.TelegramUser{}).Where("id = ?", ID).Updates(payload).Error
	if err != nil {
		return err
	}
	return nil
}

func (t *TelegramRepository) FindByUsername(username string, telegramUser *entities.TelegramUser, c *fiber.Ctx, tx *gorm.DB) error {
	err := tx.WithContext(c.Context()).Where("username = ?", username).First(&telegramUser).Error
	if err != nil {
		return err
	}
	return nil
}

func (t *TelegramRepository) DeleteByTelegramID(telegramID int64, c *fiber.Ctx, tx *gorm.DB) error {
	err := tx.WithContext(c.Context()).Where("telegram_id = ?", telegramID).Delete(&entities.TelegramUser{}).Error
	if err != nil {
		return err
	}
	return nil
}
