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

func (t *TelegramService) DeleteByTelegramID(telegramID int64, c *fiber.Ctx, tx *gorm.DB) error {
	helpers.MyLogger("debug", "TelegramAccountLink", "DeleteByTelegramId", "service", "start delete telegram user by telegram ID", nil, c)
	var telegramUser entities.TelegramUser
	// find telegram user by telegram ID
	helpers.MyLogger("debug", "TelegramAccountLink", "DeleteByTelegramId", "service", "start calling repository for find telegram user by telegram ID", nil, c)
	if err := t.TelegramRepository.FindByTelegramID(telegramID, &telegramUser, c, tx); err != nil {
		if helpers.IsNotFoundError(err) {
			helpers.MyLogger("info", "TelegramAccountLink", "DeleteByTelegramId", "service", "not found telegram user by telegram ID", nil, c)
			tx.Rollback()
			return helpers.Response(c, fiber.StatusNotFound, "Telegram user not found", nil)
		}
		helpers.MyLogger("error", "TelegramAccountLink", "DeleteByTelegramId", "service", "error find telegram user by telegram ID", map[string]interface{}{
			"error": err.Error(),
		}, c)
		tx.Rollback()
		return err
	}
	err := t.TelegramRepository.DeleteByTelegramID(telegramUser.TelegramID, c, tx)
	if err != nil {
		tx.Rollback()
		helpers.MyLogger("error", "TelegramAccountLink", "DeleteByTelegramId", "service", "error delete telegram user by telegram ID", nil, c)
		return err
	}
	if err := tx.Commit().Error; err != nil {
		helpers.MyLogger("error", "TelegramAccountLink", "DeleteByTelegramId", "service", "error committing transaction", map[string]interface{}{
			"error": err.Error(),
		}, c)
		tx.Rollback()
		return err
	}
	helpers.MyLogger("info", "TelegramAccountLink", "DeleteByTelegramId", "service", "success delete telegram user by telegram ID", nil, c)
	return helpers.Response(c, fiber.StatusOK, "Telegram user deleted successfully", nil)
}

// find by user id
func (t *TelegramService) FindByUserID(userID uint, c *fiber.Ctx, tx *gorm.DB) error {
	helpers.MyLogger("debug", "TelegramAccountLink", "FindByUserId", "service", "start find telegram user by user ID", nil, c)
	var telegramUser []entities.TelegramUser
	if err := t.TelegramRepository.FindByUserID(userID, &telegramUser, c, tx); err != nil {
		if helpers.IsNotFoundError(err) {
			helpers.MyLogger("info", "TelegramAccountLink", "FindByUserId", "service", "not found telegram user by user ID", map[string]interface{}{
				"user_id_throw": userID,
			}, c)
			// tx.Rollback()
			return helpers.Response(c, fiber.StatusNotFound, "Telegram user not found", nil)
		}
		helpers.MyLogger("error", "TelegramAccountLink", "FindByUserId", "service", "error find telegram user by user ID", map[string]interface{}{
			"error": err.Error(),
		}, c)
		return err
	}
	helpers.MyLogger("info", "TelegramAccountLink", "FindByUserId", "service", "success find telegram user by user ID", map[string]interface{}{
		"telegram_user": telegramUser,
	}, c)
	return helpers.Response(c, fiber.StatusOK, "Telegram user found successfully", telegramUser)
}

func (t *TelegramService) FindByTelegramID(telegramID int64, c *fiber.Ctx, tx *gorm.DB) error {
	helpers.MyLogger("debug", "TelegramAccountLink", "FindByTelegramId", "service", "start find telegram user by telegram ID", nil, c)
	var telegramUser entities.TelegramUser
	if err := t.TelegramRepository.FindByTelegramID(telegramID, &telegramUser, c, tx); err != nil {
		if helpers.IsNotFoundError(err) {
			helpers.MyLogger("info", "TelegramAccountLink", "FindByTelegramId", "service", "not found telegram user by telegram ID", map[string]interface{}{
				"telegram_id_throw": telegramID,
			}, c)
			return helpers.Response(c, fiber.StatusNotFound, "Telegram user not found", nil)
		}
		helpers.MyLogger("error", "TelegramAccountLink", "FindByTelegramId", "service", "error find telegram user by telegram ID", map[string]interface{}{
			"error": err.Error(),
		}, c)
		return err
	}
	helpers.MyLogger("info", "TelegramAccountLink", "FindByTelegramId", "service", "success find telegram user by telegram ID", map[string]interface{}{
		"telegram_user": telegramUser,
	}, c)
	return helpers.Response(c, fiber.StatusOK, "Telegram user found successfully", telegramUser)
}


func (t *TelegramService) UpdateByTelegramID(telegramID int64, payload *payloads.UpdateTelegramPayload, c *fiber.Ctx, tx *gorm.DB) error {
	helpers.MyLogger("debug", "TelegramAccountLink", "UpdateByTelegramId", "service", "start update telegram user by telegram ID", nil, c)
	var telegramUser entities.TelegramUser
	if err := t.TelegramRepository.FindByTelegramID(telegramID, &telegramUser, c, tx); err != nil {
		if helpers.IsNotFoundError(err) {
			helpers.MyLogger("info", "TelegramAccountLink", "UpdateByTelegramId", "service", "not found telegram user by telegram ID", map[string]interface{}{
				"telegram_id_throw": telegramID,
			}, c)
			return helpers.Response(c, fiber.StatusNotFound, "Telegram user not found", nil)
		}
		helpers.MyLogger("error", "TelegramAccountLink", "UpdateByTelegramId", "service", "error find telegram user by telegram ID", map[string]interface{}{
			"error": err.Error(),
		}, c)
		return err
	}
	telegramUser.Username = payload.Username
	telegramUser.FirstName = payload.FirstName
	telegramUser.LastName = payload.LastName
	err := t.TelegramRepository.UpdateByTelegramID(telegramID, &telegramUser, c, tx)
	if err != nil {
		helpers.MyLogger("error", "TelegramAccountLink", "UpdateByTelegramId", "service", "error update telegram user by telegram ID", map[string]interface{}{
			"error": err.Error(),
		}, c)
		tx.Rollback()
		return err
	}
	if err := tx.Commit().Error; err != nil {
		helpers.MyLogger("error", "TelegramAccountLink", "UpdateByTelegramId", "service", "error committing transaction", map[string]interface{}{
			"error": err.Error(),
		}, c)
		tx.Rollback()
		return err
	}
	helpers.MyLogger("info", "TelegramAccountLink", "UpdateByTelegramId", "service", "success update telegram user by telegram ID", map[string]interface{}{
		"telegram_user": telegramUser,
	}, c)
	return helpers.Response(c, fiber.StatusOK, "Telegram user updated successfully", telegramUser)
}