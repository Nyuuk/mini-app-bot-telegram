package repositories

import (
	"github.com/Nyuuk/mini-app-bot-telegram/backend/app/entities"
	"gorm.io/gorm"
)

type ApiKeyRepository struct{}

func (r *ApiKeyRepository) FindByApiKey(apiKey string, apiKeyEntity *entities.APIKey, tx *gorm.DB) error {
	if err := tx.Preload("User").Where("api_key = ? and is_active = ?", apiKey, true).First(&apiKeyEntity).Error; err != nil {
		return err
	}
	return nil
}

func (r *ApiKeyRepository) FindByUserID(userID uint, apiKeyEntity *entities.APIKey, tx *gorm.DB) error {
	if err := tx.Where("user_id = ?", userID).First(&apiKeyEntity).Error; err != nil {
		return err
	}
	return nil
}
