package services

import (
	"strconv"

	"github.com/Nyuuk/mini-app-bot-telegram/backend/app/entities"
	"github.com/Nyuuk/mini-app-bot-telegram/backend/app/payloads"
	"github.com/Nyuuk/mini-app-bot-telegram/backend/app/pkg/helpers"
	"github.com/Nyuuk/mini-app-bot-telegram/backend/app/repositories"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type TelegramService struct {
	TelegramRepository repositories.TelegramRepository
}

func (t *TelegramService) CreateNewUserForNowUserActive(payload *payloads.CreateNewTelegramPayload, c *fiber.Ctx, tx *gorm.DB) error {
	userID := helpers.GetCurrentUserID(c)
	helpers.LogBusiness("telegram_create_new_user_telegram", strconv.Itoa(int(userID)), map[string]interface{}{
		"user_id":     userID,
		"telegram_id": payload.TelegramID,
		"username":    payload.Username,
		"first_name":  payload.FirstName,
		"last_name":   payload.LastName,
	})
	helpers.Logger.Debug().Msg("TelegramService: CreateNewUserForNowUserActive Calling Repository for create new telegram user")
	var telegramUser entities.TelegramUser
	telegramUser.UserID = userID
	telegramUser.TelegramID = payload.TelegramID
	telegramUser.Username = payload.Username
	telegramUser.FirstName = payload.FirstName
	telegramUser.LastName = payload.LastName
	err := t.TelegramRepository.Create(&telegramUser, c, tx)
	if err != nil {
		if helpers.IsDuplicateKeyError(err) {
			helpers.Logger.Info().Msg("TelegramService: CreateNewUserForNowUserActive Telegram user already exists")
			tx.Rollback()
			return helpers.Response(c, fiber.StatusOK, "Telegram user already exists", nil)
		}
		helpers.Logger.Error().Err(err).Msg("TelegramService: CreateNewUserForNowUserActive Error creating telegram user")
		tx.Rollback()
		return err
	}
	if err := tx.Commit().Error; err != nil {
		helpers.Logger.Error().Err(err).Msg("TelegramService: CreateNewUserForNowUserActive Error committing transaction")
		tx.Rollback()
		return err
	}
	helpers.Logger.Info().Interface("telegramUser", telegramUser).Msg("TelegramService: CreateNewUserForNowUserActive Telegram user created successfully")
	return helpers.Response(c, fiber.StatusCreated, "Telegram user created successfully", telegramUser)
}
