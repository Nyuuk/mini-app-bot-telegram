package repositories

import (
	"github.com/Nyuuk/mini-app-bot-telegram/backend/app/entities"
	"gorm.io/gorm"
)

type ApiKeyRepository struct{}

func (r *ApiKeyRepository) FindByApiKey(apiKey string, apiKeyEntity *entities.APIKey, tx *gorm.DB) error {
	if err := tx.Where("api_key = ? and is_active = ?", apiKey, true).First(&apiKeyEntity).Error; err != nil {
		return err
	}
	return nil
}
