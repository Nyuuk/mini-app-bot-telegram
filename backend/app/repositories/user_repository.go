package repositories

import (
	"github.com/Nyuuk/mini-app-bot-telegram/backend/app/entities"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type UserRepository struct{}

func (r *UserRepository) FindByID(id uint, user *entities.User, tx *gorm.DB) error {
	if err := tx.Where("id = ?", id).First(&user).Error; err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) CreateUser(user *entities.User, tx *gorm.DB, c *fiber.Ctx) error {
	if err := tx.WithContext(c.Context()).Create(&user).Error; err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) FindByUsername(username string, user *entities.User, tx *gorm.DB) error {
	if err := tx.Where("username = ?", username).First(&user).Error; err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) CreateApiKey(apiKey *entities.APIKey, tx *gorm.DB, c *fiber.Ctx) error {
	if err := tx.WithContext(c.Context()).Create(&apiKey).Error; err != nil {
		return err
	}
	return nil
}