package services

import (
	"time"

	"github.com/Nyuuk/mini-app-bot-telegram/backend/app/entities"
	"github.com/Nyuuk/mini-app-bot-telegram/backend/app/payloads"
	"github.com/Nyuuk/mini-app-bot-telegram/backend/app/pkg/database"
	"github.com/Nyuuk/mini-app-bot-telegram/backend/app/pkg/helpers"
	"github.com/Nyuuk/mini-app-bot-telegram/backend/app/repositories"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type OvertimeService struct {
	OvertimeRepository repositories.OvertimeRepository
}

// CreateNewRecordOvertime creates a new overtime record
func (o *OvertimeService) CreateNewRecordOvertime(payload *payloads.CreateNewRecordOvertime, c *fiber.Ctx, tx *gorm.DB) error {
	userID := helpers.GetCurrentUserID(c)
	helpers.MyLogger("debug", "OvertimeManagement", "CreateNewRecordOvertime", "service", "start create new overtime record", map[string]interface{}{
		"user_id":        userID,
		"telegram_id":    payload.TelegramID,
		"date":           payload.Date,
		"time_start":     payload.TimeStart,
		"time_stop":      payload.TimeStop,
		"break_duration": payload.BreakDuration,
		"duration":       payload.Duration,
		"description":    payload.Description,
		"category":       payload.Category,
	}, c)

	// Get telegram_user_id from telegram_id
	helpers.MyLogger("debug", "OvertimeManagement", "CreateNewRecordOvertime", "service", "get telegram user id | calling repository Overtime GetTelegramUserIDByTelegramID", map[string]interface{}{
		"telegram_id": payload.TelegramID,
	}, c)
	telegramUserID, err := o.OvertimeRepository.GetTelegramUserIDByTelegramID(payload.TelegramID, c, tx)
	if err != nil {
		if helpers.IsNotFoundError(err) {
			helpers.MyLogger("info", "OvertimeManagement", "CreateNewRecordOvertime", "service", "telegram user not found", map[string]interface{}{
				"telegram_id": payload.TelegramID,
			}, c)
			tx.Rollback()
			return helpers.Response(c, fiber.StatusNotFound, "Telegram user not found", nil)
		}
		helpers.MyLogger("error", "OvertimeManagement", "CreateNewRecordOvertime", "service", "error finding telegram user", map[string]interface{}{
			"error": err.Error(),
		}, c)
		tx.Rollback()
		return err
	}

	// Parse datetime strings with Asia/Jakarta timezone
	helpers.MyLogger("debug", "OvertimeManagement", "CreateNewRecordOvertime", "service", "parsing date and time with timezone Asia/Jakarta | calling helpers.ParseDateWithTimezone", map[string]interface{}{
		"date":        payload.Date,
		"time_start":  payload.TimeStart,
		"time_stop":   payload.TimeStop,
		"telegram_id": payload.TelegramID,
	}, c)
	date, err := helpers.ParseDateWithTimezone(payload.Date)
	if err != nil {
		helpers.MyLogger("error", "OvertimeManagement", "CreateNewRecordOvertime", "service", "error parsing date", map[string]interface{}{
			"error": err.Error(),
			"date":  payload.Date,
		}, c)
		tx.Rollback()
		return helpers.Response(c, fiber.StatusBadRequest, "Invalid date format. Use YYYY-MM-DD", nil)
	}
	helpers.MyLogger("debug", "OvertimeManagement", "CreateNewRecordOvertime", "service", "successfully parsed date", map[string]interface{}{
		"date": date,
	}, c)

	// Parse time_start - use appropriate parser based on format
	var timeStart time.Time
	if helpers.IsTimeOnlyFormat(payload.TimeStart) {
		helpers.MyLogger("debug", "OvertimeManagement", "CreateNewRecordOvertime", "service", "calling helpers.ParseTimeWithTimezone", map[string]interface{}{
			"time_start": payload.TimeStart,
		}, c)
		timeStart, err = helpers.ParseTimeWithTimezone(payload.TimeStart)
	} else {
		helpers.MyLogger("debug", "OvertimeManagement", "CreateNewRecordOvertime", "service", "calling helpers.ParseDateTimeWithTimezone", map[string]interface{}{
			"time_start": payload.TimeStart,
		}, c)
		timeStart, err = helpers.ParseDateTimeWithTimezone(payload.TimeStart)
	}
	if err != nil {
		helpers.MyLogger("error", "OvertimeManagement", "CreateNewRecordOvertime", "service", "error parsing time start", map[string]interface{}{
			"error":      err.Error(),
			"time_start": payload.TimeStart,
		}, c)
		tx.Rollback()
		return helpers.Response(c, fiber.StatusBadRequest, "Invalid time start format. Use YYYY-MM-DDTHH:MM:SS or HH:MM:SS", nil)
	}
	helpers.MyLogger("debug", "OvertimeManagement", "CreateNewRecordOvertime", "service", "successfully parsed time start", map[string]interface{}{
		"time_start": timeStart,
	}, c)

	// Parse time_stop - use appropriate parser based on format
	var timeStop time.Time
	if helpers.IsTimeOnlyFormat(payload.TimeStop) {
		helpers.MyLogger("debug", "OvertimeManagement", "CreateNewRecordOvertime", "service", "calling helpers.ParseTimeWithTimezone", map[string]interface{}{
			"time_stop": payload.TimeStop,
		}, c)
		timeStop, err = helpers.ParseTimeWithTimezone(payload.TimeStop)
	} else {
		helpers.MyLogger("debug", "OvertimeManagement", "CreateNewRecordOvertime", "service", "calling helpers.ParseDateTimeWithTimezone", map[string]interface{}{
			"time_stop": payload.TimeStop,
		}, c)
		timeStop, err = helpers.ParseDateTimeWithTimezone(payload.TimeStop)
	}
	if err != nil {
		helpers.MyLogger("error", "OvertimeManagement", "CreateNewRecordOvertime", "service", "error parsing time stop", map[string]interface{}{
			"error":     err.Error(),
			"time_stop": payload.TimeStop,
		}, c)
		tx.Rollback()
		return helpers.Response(c, fiber.StatusBadRequest, "Invalid time stop format. Use YYYY-MM-DDTHH:MM:SS or HH:MM:SS", nil)
	}
	helpers.MyLogger("debug", "OvertimeManagement", "CreateNewRecordOvertime", "service", "successfully parsed time stop", map[string]interface{}{
		"time_stop": timeStop,
	}, c)

	var overtime entities.Overtime
	overtime.TelegramUserID = telegramUserID
	overtime.Date = date
	overtime.TimeStart = timeStart
	overtime.TimeStop = timeStop
	overtime.BreakDuration = payload.BreakDuration
	overtime.Duration = payload.Duration
	overtime.Description = payload.Description
	overtime.Category = payload.Category
	overtime.CreatedByUserID = userID

	helpers.MyLogger("debug", "OvertimeManagement", "CreateNewRecordOvertime", "service", "calling repository to create overtime record", nil, c)
	err = o.OvertimeRepository.CreateNewRecordOvertime(&overtime, c, tx)
	if err != nil {
		if helpers.IsDuplicateKeyError(err) {
			helpers.MyLogger("info", "OvertimeManagement", "CreateNewRecordOvertime", "service", "overtime record already exists for this date", map[string]interface{}{
				"telegram_id": payload.TelegramID,
				"date":        payload.Date,
			}, c)
			tx.Rollback()
			return helpers.Response(c, fiber.StatusConflict, "Overtime record already exists for this date", nil)
		}
		helpers.MyLogger("error", "OvertimeManagement", "CreateNewRecordOvertime", "service", "error creating overtime record", map[string]interface{}{
			"error": err.Error(),
		}, c)
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		helpers.MyLogger("error", "OvertimeManagement", "CreateNewRecordOvertime", "service", "error committing transaction", map[string]interface{}{
			"error": err.Error(),
		}, c)
		tx.Rollback()
		return err
	}

	helpers.MyLogger("info", "OvertimeManagement", "CreateNewRecordOvertime", "service", "overtime record created successfully", map[string]interface{}{
		"overtime_id": overtime.ID,
	}, c)
	return helpers.Response(c, fiber.StatusCreated, "Overtime record created successfully", overtime)
}

// GetAllRecordOvertimeByTelegramID retrieves all overtime records by telegram ID
func (o *OvertimeService) GetAllRecordOvertimeByTelegramID(telegramID int64, c *fiber.Ctx, tx *gorm.DB) error {
	helpers.MyLogger("debug", "OvertimeManagement", "GetAllRecordOvertimeByTelegramID", "service", "start get all overtime records by telegram ID", map[string]interface{}{
		"telegram_id": telegramID,
	}, c)

	var overtimes []entities.Overtime
	err := o.OvertimeRepository.GetAllRecordOvertimeByTelegramID(telegramID, &overtimes, c, tx)
	if err != nil {
		if helpers.IsNotFoundError(err) {
			helpers.MyLogger("info", "OvertimeManagement", "GetAllRecordOvertimeByTelegramID", "service", "no overtime records found for telegram user", map[string]interface{}{
				"telegram_id": telegramID,
			}, c)
			return helpers.Response(c, fiber.StatusNotFound, "No overtime records found", nil)
		}
		helpers.MyLogger("error", "OvertimeManagement", "GetAllRecordOvertimeByTelegramID", "service", "error getting overtime records", map[string]interface{}{
			"error": err.Error(),
		}, c)
		return err
	}

	helpers.MyLogger("info", "OvertimeManagement", "GetAllRecordOvertimeByTelegramID", "service", "overtime records retrieved successfully", map[string]interface{}{
		"telegram_id":   telegramID,
		"records_count": len(overtimes),
	}, c)
	return helpers.Response(c, fiber.StatusOK, "Overtime records retrieved successfully", overtimes)
}

// GetRecordByDateByTelegramID retrieves overtime record by specific date and telegram ID
func (o *OvertimeService) GetRecordByDateByTelegramID(telegramID int64, date time.Time, c *fiber.Ctx, tx *gorm.DB) error {
	helpers.MyLogger("debug", "OvertimeManagement", "GetRecordByDateByTelegramID", "service", "start get overtime record by date and telegram ID", map[string]interface{}{
		"telegram_id": telegramID,
		"date":        date,
	}, c)

	var overtime []entities.Overtime
	err := o.OvertimeRepository.GetRecordByDateByTelegramID(telegramID, date, &overtime, c, tx)
	if err != nil {
		if helpers.IsNotFoundError(err) {
			helpers.MyLogger("info", "OvertimeManagement", "GetRecordByDateByTelegramID", "service", "overtime record not found for date", map[string]interface{}{
				"telegram_id": telegramID,
				"date":        date,
			}, c)
			return helpers.Response(c, fiber.StatusNotFound, "Overtime record not found for this date", nil)
		}
		helpers.MyLogger("error", "OvertimeManagement", "GetRecordByDateByTelegramID", "service", "error getting overtime record by date", map[string]interface{}{
			"error": err.Error(),
		}, c)
		return err
	}

	helpers.MyLogger("info", "OvertimeManagement", "GetRecordByDateByTelegramID", "service", "overtime record retrieved successfully", map[string]interface{}{
		"telegram_id": telegramID,
		"date":        date,
		"overtime":    overtime,
	}, c)
	return helpers.Response(c, fiber.StatusOK, "Overtime record retrieved successfully", overtime)
}

// GetRecordBetweenDateByTelegramId retrieves overtime records between two dates for a telegram ID
func (o *OvertimeService) GetRecordBetweenDateByTelegramId(telegramID int64, startDate time.Time, endDate time.Time, c *fiber.Ctx, tx *gorm.DB) error {
	helpers.MyLogger("debug", "OvertimeManagement", "GetRecordBetweenDateByTelegramId", "service", "start get overtime records between dates", map[string]interface{}{
		"telegram_id": telegramID,
		"start_date":  startDate,
		"end_date":    endDate,
	}, c)

	var overtimes []entities.Overtime
	err := o.OvertimeRepository.GetRecordBetweenDateByTelegramId(telegramID, startDate, endDate, &overtimes, c, tx)
	if err != nil {
		if helpers.IsNotFoundError(err) {
			helpers.MyLogger("info", "OvertimeManagement", "GetRecordBetweenDateByTelegramId", "service", "no overtime records found between dates", map[string]interface{}{
				"telegram_id": telegramID,
				"start_date":  startDate,
				"end_date":    endDate,
			}, c)
			return helpers.Response(c, fiber.StatusNotFound, "No overtime records found between these dates", nil)
		}
		helpers.MyLogger("error", "OvertimeManagement", "GetRecordBetweenDateByTelegramId", "service", "error getting overtime records between dates", map[string]interface{}{
			"error": err.Error(),
		}, c)
		return err
	}

	// Calculate total duration for the period
	var totalDuration float64
	for _, overtime := range overtimes {
		totalDuration += overtime.Duration
	}

	helpers.MyLogger("info", "OvertimeManagement", "GetRecordBetweenDateByTelegramId", "service", "overtime records retrieved successfully", map[string]interface{}{
		"telegram_id":    telegramID,
		"start_date":     startDate,
		"end_date":       endDate,
		"records_count":  len(overtimes),
		"total_duration": totalDuration,
	}, c)

	response := map[string]interface{}{
		"records":        overtimes,
		"total_duration": totalDuration,
		"records_count":  len(overtimes),
		"period": map[string]interface{}{
			"start_date": startDate,
			"end_date":   endDate,
		},
	}

	return helpers.Response(c, fiber.StatusOK, "Overtime records retrieved successfully", response)
}

// GetRecordByID retrieves overtime record by ID
func (o *OvertimeService) GetRecordByID(id uint, c *fiber.Ctx, tx *gorm.DB) error {
	helpers.MyLogger("debug", "OvertimeManagement", "GetRecordByID", "service", "start get overtime record by ID", map[string]interface{}{
		"overtime_id": id,
	}, c)

	var overtime entities.Overtime
	err := o.OvertimeRepository.GetRecordByID(id, &overtime, c, tx)
	if err != nil {
		if helpers.IsNotFoundError(err) {
			helpers.MyLogger("info", "OvertimeManagement", "GetRecordByID", "service", "overtime record not found", map[string]interface{}{
				"overtime_id": id,
			}, c)
			return helpers.Response(c, fiber.StatusNotFound, "Overtime record not found", nil)
		}
		helpers.MyLogger("error", "OvertimeManagement", "GetRecordByID", "service", "error getting overtime record by ID", map[string]interface{}{
			"error": err.Error(),
		}, c)
		return err
	}

	helpers.MyLogger("info", "OvertimeManagement", "GetRecordByID", "service", "overtime record retrieved successfully", map[string]interface{}{
		"overtime_id": id,
	}, c)
	return helpers.Response(c, fiber.StatusOK, "Overtime record retrieved successfully", overtime)
}

// UpdateRecordOvertime updates an existing overtime record
func (o *OvertimeService) UpdateRecordOvertime(id uint, payload *payloads.UpdateRecordOvertime, c *fiber.Ctx, tx *gorm.DB) error {
	userID := helpers.GetCurrentUserID(c)
	helpers.MyLogger("debug", "OvertimeManagement", "UpdateRecordOvertime", "service", "start update overtime record", map[string]interface{}{
		"user_id":        userID,
		"overtime_id":    id,
		"telegram_id":    payload.TelegramID,
		"date":           payload.Date,
		"time_start":     payload.TimeStart,
		"time_stop":      payload.TimeStop,
		"duration":       payload.Duration,
		"break_duration": payload.BreakDuration,
		"description":    payload.Description,
		"category":       payload.Category,
	}, c)

	// Check if record exists
	helpers.MyLogger("debug", "OvertimeManagement", "UpdateRecordOvertime", "service", "check if record exists | calling repository GetRecordByID", map[string]interface{}{
		"overtime_id": id,
	}, c)
	var existingOvertime entities.Overtime
	err := o.OvertimeRepository.GetRecordByID(id, &existingOvertime, c, tx)
	if err != nil {
		if helpers.IsNotFoundError(err) {
			helpers.MyLogger("info", "OvertimeManagement", "UpdateRecordOvertime", "service", "overtime record not found", map[string]interface{}{
				"overtime_id": id,
			}, c)
			return helpers.Response(c, fiber.StatusNotFound, "Overtime record not found", nil)
		}
		helpers.MyLogger("error", "OvertimeManagement", "UpdateRecordOvertime", "service", "error finding overtime record", map[string]interface{}{
			"error": err.Error(),
		}, c)
		return err
	}

	// Prepare update data - only update fields that are provided
	updates := make(map[string]interface{})

	// Handle TelegramID update
	if payload.TelegramID != 0 {
		helpers.MyLogger("debug", "OvertimeManagement", "UpdateRecordOvertime", "service", "get telegram user id | calling repository GetTelegramUserIDByTelegramID", map[string]interface{}{
			"telegram_id": payload.TelegramID,
		}, c)
		telegramUserID, err := o.OvertimeRepository.GetTelegramUserIDByTelegramID(payload.TelegramID, c, tx)
		if err != nil {
			if helpers.IsNotFoundError(err) {
				helpers.MyLogger("info", "OvertimeManagement", "UpdateRecordOvertime", "service", "telegram user not found", map[string]interface{}{
					"telegram_id": payload.TelegramID,
				}, c)
				tx.Rollback()
				return helpers.Response(c, fiber.StatusNotFound, "Telegram user not found", nil)
			}
			helpers.MyLogger("error", "OvertimeManagement", "UpdateRecordOvertime", "service", "error finding telegram user", map[string]interface{}{
				"error": err.Error(),
			}, c)
			tx.Rollback()
			return err
		}
		updates["telegram_user_id"] = telegramUserID
		helpers.MyLogger("debug", "OvertimeManagement", "UpdateRecordOvertime", "service", "telegram_user_id will be updated", map[string]interface{}{
			"telegram_user_id": telegramUserID,
		}, c)
	}

	// Handle Date update
	if payload.Date != "" {
		helpers.MyLogger("debug", "OvertimeManagement", "UpdateRecordOvertime", "service", "parsing date with timezone Asia/Jakarta | calling helpers.ParseDateWithTimezone", map[string]interface{}{
			"date": payload.Date,
		}, c)
		date, err := helpers.ParseDateWithTimezone(payload.Date)
		if err != nil {
			helpers.MyLogger("error", "OvertimeManagement", "UpdateRecordOvertime", "service", "error parsing date", map[string]interface{}{
				"error": err.Error(),
				"date":  payload.Date,
			}, c)
			tx.Rollback()
			return helpers.Response(c, fiber.StatusBadRequest, "Invalid date format. Use YYYY-MM-DD", nil)
		}
		updates["date"] = date
		helpers.MyLogger("debug", "OvertimeManagement", "UpdateRecordOvertime", "service", "date will be updated", map[string]interface{}{
			"date": date,
		}, c)
	}

	// Handle TimeStart update
	if payload.TimeStart != "" {
		var timeStart time.Time
		if helpers.IsTimeOnlyFormat(payload.TimeStart) {
			helpers.MyLogger("debug", "OvertimeManagement", "UpdateRecordOvertime", "service", "calling helpers.ParseTimeWithTimezone", map[string]interface{}{
				"time_start": payload.TimeStart,
			}, c)
			timeStart, err = helpers.ParseTimeWithTimezone(payload.TimeStart)
		} else {
			helpers.MyLogger("debug", "OvertimeManagement", "UpdateRecordOvertime", "service", "calling helpers.ParseDateTimeWithTimezone", map[string]interface{}{
				"time_start": payload.TimeStart,
			}, c)
			timeStart, err = helpers.ParseDateTimeWithTimezone(payload.TimeStart)
		}
		if err != nil {
			helpers.MyLogger("error", "OvertimeManagement", "UpdateRecordOvertime", "service", "error parsing time start", map[string]interface{}{
				"error":      err.Error(),
				"time_start": payload.TimeStart,
			}, c)
			tx.Rollback()
			return helpers.Response(c, fiber.StatusBadRequest, "Invalid time start format. Use YYYY-MM-DDTHH:MM:SS or HH:MM:SS", nil)
		}
		updates["time_start"] = timeStart
		helpers.MyLogger("debug", "OvertimeManagement", "UpdateRecordOvertime", "service", "time_start will be updated", map[string]interface{}{
			"time_start": timeStart,
		}, c)
	}

	// Handle TimeStop update
	if payload.TimeStop != "" {
		var timeStop time.Time
		if helpers.IsTimeOnlyFormat(payload.TimeStop) {
			helpers.MyLogger("debug", "OvertimeManagement", "UpdateRecordOvertime", "service", "calling helpers.ParseTimeWithTimezone", map[string]interface{}{
				"time_stop": payload.TimeStop,
			}, c)
			timeStop, err = helpers.ParseTimeWithTimezone(payload.TimeStop)
		} else {
			helpers.MyLogger("debug", "OvertimeManagement", "UpdateRecordOvertime", "service", "calling helpers.ParseDateTimeWithTimezone", map[string]interface{}{
				"time_stop": payload.TimeStop,
			}, c)
			timeStop, err = helpers.ParseDateTimeWithTimezone(payload.TimeStop)
		}
		if err != nil {
			helpers.MyLogger("error", "OvertimeManagement", "UpdateRecordOvertime", "service", "error parsing time stop", map[string]interface{}{
				"error":     err.Error(),
				"time_stop": payload.TimeStop,
			}, c)
			tx.Rollback()
			return helpers.Response(c, fiber.StatusBadRequest, "Invalid time stop format. Use YYYY-MM-DDTHH:MM:SS or HH:MM:SS", nil)
		}
		updates["time_stop"] = timeStop
		helpers.MyLogger("debug", "OvertimeManagement", "UpdateRecordOvertime", "service", "time_stop will be updated", map[string]interface{}{
			"time_stop": timeStop,
		}, c)
	}

	// Handle Description update
	if payload.Description != "" {
		updates["description"] = payload.Description
		helpers.MyLogger("debug", "OvertimeManagement", "UpdateRecordOvertime", "service", "description will be updated", map[string]interface{}{
			"description": payload.Description,
		}, c)
	}

	// Handle Category update
	if payload.Category != "" {
		updates["category"] = payload.Category
		helpers.MyLogger("debug", "OvertimeManagement", "UpdateRecordOvertime", "service", "category will be updated", map[string]interface{}{
			"category": payload.Category,
		}, c)
	}

	// Handle Duration update
	if payload.Duration > 0 {
		updates["duration"] = payload.Duration
		helpers.MyLogger("debug", "OvertimeManagement", "UpdateRecordOvertime", "service", "duration will be updated", map[string]interface{}{
			"duration": payload.Duration,
		}, c)
	}

	// Handle BreakDuration update
	if payload.BreakDuration >= 0 {
		updates["break_duration"] = payload.BreakDuration
		helpers.MyLogger("debug", "OvertimeManagement", "UpdateRecordOvertime", "service", "break_duration will be updated", map[string]interface{}{
			"break_duration": payload.BreakDuration,
		}, c)
	}

	// Check if any fields are being updated
	if len(updates) == 0 {
		helpers.MyLogger("info", "OvertimeManagement", "UpdateRecordOvertime", "service", "no fields to update", map[string]interface{}{
			"overtime_id": id,
		}, c)
		tx.Rollback()
		return helpers.Response(c, fiber.StatusBadRequest, "No fields to update", nil)
	}

	helpers.MyLogger("debug", "OvertimeManagement", "UpdateRecordOvertime", "service", "fields to be updated", map[string]interface{}{
		"updates": updates,
	}, c)

	helpers.MyLogger("debug", "OvertimeManagement", "UpdateRecordOvertime", "service", "calling repository to update overtime record", map[string]interface{}{
		"overtime_id": id,
	}, c)
	err = o.OvertimeRepository.UpdateRecordOvertimePartial(id, updates, c, tx)
	if err != nil {
		helpers.MyLogger("error", "OvertimeManagement", "UpdateRecordOvertime", "service", "error updating overtime record", map[string]interface{}{
			"error": err.Error(),
		}, c)
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		helpers.MyLogger("error", "OvertimeManagement", "UpdateRecordOvertime", "service", "error committing transaction", map[string]interface{}{
			"error": err.Error(),
		}, c)
		tx.Rollback()
		return err
	}

	// Get updated record to return to client
	helpers.MyLogger("debug", "OvertimeManagement", "UpdateRecordOvertime", "service", "get updated record | calling repository GetRecordByID", map[string]interface{}{
		"overtime_id": id,
	}, c)
	var updatedRecord entities.Overtime
	err = o.OvertimeRepository.GetRecordByID(id, &updatedRecord, c, database.ClientPostgres)
	if err != nil {
		helpers.MyLogger("error", "OvertimeManagement", "UpdateRecordOvertime", "service", "error getting updated record", map[string]interface{}{
			"error": err.Error(),
		}, c)
		return helpers.Response(c, fiber.StatusInternalServerError, "Record updated but failed to retrieve updated data", nil)
	}

	helpers.MyLogger("info", "OvertimeManagement", "UpdateRecordOvertime", "service", "overtime record updated successfully", map[string]interface{}{
		"overtime_id":    id,
		"updated_fields": len(updates),
	}, c)
	return helpers.Response(c, fiber.StatusOK, "Overtime record updated successfully", updatedRecord)
}

// DeleteRecordOvertime deletes an overtime record
func (o *OvertimeService) DeleteRecordOvertime(id uint, c *fiber.Ctx, tx *gorm.DB) error {
	userID := helpers.GetCurrentUserID(c)
	helpers.MyLogger("debug", "OvertimeManagement", "DeleteRecordOvertime", "service", "start delete overtime record", map[string]interface{}{
		"overtime_id": id,
		"user_id":     userID,
	}, c)

	// Check if record exists
	helpers.MyLogger("debug", "OvertimeManagement", "DeleteRecordOvertime", "service", "check if record exists | calling repository GetRecordByID", map[string]interface{}{
		"overtime_id": id,
	}, c)
	var overtime entities.Overtime
	err := o.OvertimeRepository.GetRecordByID(id, &overtime, c, tx)
	if err != nil {
		if helpers.IsNotFoundError(err) {
			helpers.MyLogger("info", "OvertimeManagement", "DeleteRecordOvertime", "service", "overtime record not found", map[string]interface{}{
				"overtime_id": id,
			}, c)
			return helpers.Response(c, fiber.StatusNotFound, "Overtime record not found", nil)
		}
		helpers.MyLogger("error", "OvertimeManagement", "DeleteRecordOvertime", "service", "error finding overtime record", map[string]interface{}{
			"error": err.Error(),
		}, c)
		return err
	}

	helpers.MyLogger("debug", "OvertimeManagement", "DeleteRecordOvertime", "service", "record found, proceeding with deletion", map[string]interface{}{
		"overtime_id":      id,
		"description":      overtime.Description,
		"telegram_user_id": overtime.TelegramUserID,
	}, c)

	helpers.MyLogger("debug", "OvertimeManagement", "DeleteRecordOvertime", "service", "calling repository to delete overtime record", map[string]interface{}{
		"overtime_id": id,
	}, c)
	err = o.OvertimeRepository.DeleteRecordOvertime(id, c, tx)
	if err != nil {
		helpers.MyLogger("error", "OvertimeManagement", "DeleteRecordOvertime", "service", "error deleting overtime record", map[string]interface{}{
			"error": err.Error(),
		}, c)
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		helpers.MyLogger("error", "OvertimeManagement", "DeleteRecordOvertime", "service", "error committing transaction", map[string]interface{}{
			"error": err.Error(),
		}, c)
		tx.Rollback()
		return err
	}

	helpers.MyLogger("info", "OvertimeManagement", "DeleteRecordOvertime", "service", "overtime record deleted successfully", map[string]interface{}{
		"overtime_id": id,
		"deleted_by":  userID,
	}, c)
	return helpers.Response(c, fiber.StatusOK, "Overtime record deleted successfully", nil)
}
