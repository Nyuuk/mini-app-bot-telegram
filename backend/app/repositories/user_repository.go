package repositories

import (
	"time"

	"github.com/Nyuuk/mini-app-bot-telegram/backend/app/entities"
	"github.com/Nyuuk/mini-app-bot-telegram/backend/app/pkg/helpers"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type UserRepository struct{}

func (r *UserRepository) FindByID(id uint, user *entities.User, tx *gorm.DB) error {
	if err := tx.Preload("APIKeys").Where("id = ?", id).First(&user).Error; err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) CreateUser(user *entities.User, tx *gorm.DB, c *fiber.Ctx) error {
	start := time.Now()

	if err := tx.WithContext(c.Context()).Create(&user).Error; err != nil {
		helpers.LogDatabase("create", user.TableName(), time.Since(start), 0, err)
		return err
	}

	helpers.LogDatabase("create", user.TableName(), time.Since(start), 1, nil)
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

func (r *UserRepository) FindAll(users *[]entities.User, tx *gorm.DB) error {
	if err := tx.Find(&users).Error; err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) DeleteUserById(id uint, tx *gorm.DB) error {
	if err := tx.Delete(&entities.User{}, id).Error; err != nil {
		return err
	}
	return nil
}
