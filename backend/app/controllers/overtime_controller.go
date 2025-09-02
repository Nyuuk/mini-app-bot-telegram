package controllers

import (
	"strconv"

	"github.com/Nyuuk/mini-app-bot-telegram/backend/app/payloads"
	"github.com/Nyuuk/mini-app-bot-telegram/backend/app/pkg/database"
	"github.com/Nyuuk/mini-app-bot-telegram/backend/app/pkg/helpers"
	"github.com/Nyuuk/mini-app-bot-telegram/backend/app/services"
	"github.com/gofiber/fiber/v2"
)

type OvertimeController struct {
	OvertimeService services.OvertimeService
}

// CreateNewRecordOvertime godoc
// @Summary Create New Overtime Record
// @Description Create a new overtime record for a telegram user
// @Tags Overtime
// @Accept json
// @Produce json
// @Param createOvertimePayload body payloads.CreateNewRecordOvertime true "Overtime record data"
// @Success 201 {object} map[string]interface{} "Overtime record created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security ApiKeyAuth
// @Security BearerAuth
// @Router /v1/overtime/ [post]
func (o *OvertimeController) CreateNewRecordOvertime(c *fiber.Ctx) error {
	helpers.MyLogger("debug", "OvertimeManagement", "CreateNewRecordOvertime", "controller", "start create new overtime record", nil, c)

	var payload payloads.CreateNewRecordOvertime
	if err := helpers.ValidateBody(&payload, c); err != nil {
		helpers.MyLogger("error", "OvertimeManagement", "CreateNewRecordOvertime", "controller", "error validate body", map[string]interface{}{
			"error": err.Error(),
		}, c)
		if customErr, ok := err.(helpers.Error); ok {
			return helpers.ResponseErrorBadRequest(c, customErr.Message, customErr.Data)
		}
		return helpers.ResponseErrorBadRequest(c, "Invalid payload", nil)
	}

	tx := database.ClientPostgres.Begin()
	defer tx.Rollback()

	helpers.MyLogger("debug", "OvertimeManagement", "CreateNewRecordOvertime", "controller", "start calling service for create overtime record", nil, c)
	if err := o.OvertimeService.CreateNewRecordOvertime(&payload, c, tx); err != nil {
		helpers.MyLogger("error", "OvertimeManagement", "CreateNewRecordOvertime", "controller", "error create overtime record", map[string]interface{}{
			"error": err.Error(),
		}, c)
		return helpers.ResponseErrorInternal(c, err)
	}
	return nil
}

// GetAllRecordOvertimeByTelegramID retrieves all overtime records by telegram ID
func (o *OvertimeController) GetAllRecordOvertimeByTelegramID(c *fiber.Ctx) error {
	helpers.MyLogger("debug", "OvertimeManagement", "GetAllRecordOvertimeByTelegramID", "controller", "start get all overtime records by telegram ID", nil, c)

	telegramIDStr := c.Params("telegram_id")
	telegramID, err := strconv.ParseInt(telegramIDStr, 10, 64)
	if err != nil {
		helpers.MyLogger("error", "OvertimeManagement", "GetAllRecordOvertimeByTelegramID", "controller", "error parse telegram ID", map[string]interface{}{
			"error": err.Error(),
		}, c)
		return helpers.ResponseErrorBadRequest(c, "Invalid telegram ID", nil)
	}

	tx := database.ClientPostgres

	helpers.MyLogger("debug", "OvertimeManagement", "GetAllRecordOvertimeByTelegramID", "controller", "start calling service for get all overtime records", map[string]interface{}{
		"telegram_id": telegramID,
	}, c)
	if err := o.OvertimeService.GetAllRecordOvertimeByTelegramID(telegramID, c, tx); err != nil {
		helpers.MyLogger("error", "OvertimeManagement", "GetAllRecordOvertimeByTelegramID", "controller", "error get all overtime records", map[string]interface{}{
			"error": err.Error(),
		}, c)
		return helpers.ResponseErrorInternal(c, err)
	}
	return nil
}

// GetRecordByDateByTelegramID godoc
// @Summary Get Overtime Record by Date
// @Description Get overtime record for a specific date and telegram ID
// @Tags Overtime
// @Accept json
// @Produce json
// @Param getRecordByDatePayload body payloads.GetRecordByDateRequest true "Date and telegram ID"
// @Success 200 {object} map[string]interface{} "Overtime record retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 404 {object} map[string]interface{} "Overtime record not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security ApiKeyAuth
// @Security BearerAuth
// @Router /v1/overtime/by-date [post]
func (o *OvertimeController) GetRecordByDateByTelegramID(c *fiber.Ctx) error {
	helpers.MyLogger("debug", "OvertimeManagement", "GetRecordByDateByTelegramID", "controller", "start get overtime record by date", nil, c)

	var payload payloads.GetRecordByDateRequest
	if err := helpers.ValidateBody(&payload, c); err != nil {
		helpers.MyLogger("error", "OvertimeManagement", "GetRecordByDateByTelegramID", "controller", "error validate body", map[string]interface{}{
			"error": err.Error(),
		}, c)
		if customErr, ok := err.(helpers.Error); ok {
			return helpers.ResponseErrorBadRequest(c, customErr.Message, customErr.Data)
		}
		return helpers.ResponseErrorBadRequest(c, "Invalid payload", nil)
	}

	// Parse date string to time.Time with timezone
	date, err := helpers.ParseDateWithTimezone(payload.Date)
	if err != nil {
		helpers.MyLogger("error", "OvertimeManagement", "GetRecordByDateByTelegramID", "controller", "error parse date", map[string]interface{}{
			"error": err.Error(),
			"date":  payload.Date,
		}, c)
		return helpers.ResponseErrorBadRequest(c, "Invalid date format. Use YYYY-MM-DD", nil)
	}

	tx := database.ClientPostgres

	helpers.MyLogger("debug", "OvertimeManagement", "GetRecordByDateByTelegramID", "controller", "start calling service for get overtime record by date", map[string]interface{}{
		"telegram_id": payload.TelegramID,
		"date":        date,
	}, c)
	if err := o.OvertimeService.GetRecordByDateByTelegramID(payload.TelegramID, date, c, tx); err != nil {
		helpers.MyLogger("error", "OvertimeManagement", "GetRecordByDateByTelegramID", "controller", "error get overtime record by date", map[string]interface{}{
			"error": err.Error(),
		}, c)
		return helpers.ResponseErrorInternal(c, err)
	}
	return nil
}

// GetRecordBetweenDateByTelegramId retrieves overtime records between two dates for a telegram user
func (o *OvertimeController) GetRecordBetweenDateByTelegramId(c *fiber.Ctx) error {
	helpers.MyLogger("debug", "OvertimeManagement", "GetRecordBetweenDateByTelegramId", "controller", "start get overtime records between dates", nil, c)

	var payload payloads.GetRecordBetweenDateRequest
	if err := helpers.ValidateBody(&payload, c); err != nil {
		helpers.MyLogger("error", "OvertimeManagement", "GetRecordBetweenDateByTelegramId", "controller", "error validate body", map[string]interface{}{
			"error": err.Error(),
		}, c)
		if customErr, ok := err.(helpers.Error); ok {
			return helpers.ResponseErrorBadRequest(c, customErr.Message, customErr.Data)
		}
		return helpers.ResponseErrorBadRequest(c, "Invalid payload", nil)
	}

	// Parse date strings to time.Time with timezone
	startDate, err := helpers.ParseDateWithTimezone(payload.StartDate)
	if err != nil {
		helpers.MyLogger("error", "OvertimeManagement", "GetRecordBetweenDateByTelegramId", "controller", "error parse start date", map[string]interface{}{
			"error":      err.Error(),
			"start_date": payload.StartDate,
		}, c)
		return helpers.ResponseErrorBadRequest(c, "Invalid start date format. Use YYYY-MM-DD", nil)
	}

	endDate, err := helpers.ParseDateWithTimezone(payload.EndDate)
	if err != nil {
		helpers.MyLogger("error", "OvertimeManagement", "GetRecordBetweenDateByTelegramId", "controller", "error parse end date", map[string]interface{}{
			"error":    err.Error(),
			"end_date": payload.EndDate,
		}, c)
		return helpers.ResponseErrorBadRequest(c, "Invalid end date format. Use YYYY-MM-DD", nil)
	}

	// Validate date range
	if endDate.Before(startDate) {
		helpers.MyLogger("error", "OvertimeManagement", "GetRecordBetweenDateByTelegramId", "controller", "end date is before start date", map[string]interface{}{
			"start_date": startDate,
			"end_date":   endDate,
		}, c)
		return helpers.ResponseErrorBadRequest(c, "End date must be after start date", nil)
	}

	tx := database.ClientPostgres

	helpers.MyLogger("debug", "OvertimeManagement", "GetRecordBetweenDateByTelegramId", "controller", "start calling service for get overtime records between dates", map[string]interface{}{
		"telegram_id": payload.TelegramID,
		"start_date":  startDate,
		"end_date":    endDate,
	}, c)
	if err := o.OvertimeService.GetRecordBetweenDateByTelegramId(payload.TelegramID, startDate, endDate, c, tx); err != nil {
		helpers.MyLogger("error", "OvertimeManagement", "GetRecordBetweenDateByTelegramId", "controller", "error get overtime records between dates", map[string]interface{}{
			"error": err.Error(),
		}, c)
		return helpers.ResponseErrorInternal(c, err)
	}
	return nil
}

// GetRecordByID retrieves overtime record by ID
func (o *OvertimeController) GetRecordByID(c *fiber.Ctx) error {
	helpers.MyLogger("debug", "OvertimeManagement", "GetRecordByID", "controller", "start get overtime record by ID", nil, c)

	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		helpers.MyLogger("error", "OvertimeManagement", "GetRecordByID", "controller", "error parse overtime ID", map[string]interface{}{
			"error": err.Error(),
		}, c)
		return helpers.ResponseErrorBadRequest(c, "Invalid overtime ID", nil)
	}

	tx := database.ClientPostgres

	helpers.MyLogger("debug", "OvertimeManagement", "GetRecordByID", "controller", "start calling service for get overtime record by ID", map[string]interface{}{
		"overtime_id": id,
	}, c)
	if err := o.OvertimeService.GetRecordByID(uint(id), c, tx); err != nil {
		helpers.MyLogger("error", "OvertimeManagement", "GetRecordByID", "controller", "error get overtime record by ID", map[string]interface{}{
			"error": err.Error(),
		}, c)
		return helpers.ResponseErrorInternal(c, err)
	}
	return nil
}

// UpdateRecordOvertime updates an existing overtime record
func (o *OvertimeController) UpdateRecordOvertime(c *fiber.Ctx) error {
	helpers.MyLogger("debug", "OvertimeManagement", "UpdateRecordOvertime", "controller", "start update overtime record", nil, c)

	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		helpers.MyLogger("error", "OvertimeManagement", "UpdateRecordOvertime", "controller", "error parse overtime ID", map[string]interface{}{
			"error": err.Error(),
		}, c)
		return helpers.ResponseErrorBadRequest(c, "Invalid overtime ID", nil)
	}

	var payload payloads.CreateNewRecordOvertime
	if err := helpers.ValidateBody(&payload, c); err != nil {
		helpers.MyLogger("error", "OvertimeManagement", "UpdateRecordOvertime", "controller", "error validate body", map[string]interface{}{
			"error": err.Error(),
		}, c)
		if customErr, ok := err.(helpers.Error); ok {
			return helpers.ResponseErrorBadRequest(c, customErr.Message, customErr.Data)
		}
		return helpers.ResponseErrorBadRequest(c, "Invalid payload", nil)
	}

	tx := database.ClientPostgres.Begin()
	defer tx.Rollback()

	helpers.MyLogger("debug", "OvertimeManagement", "UpdateRecordOvertime", "controller", "start calling service for update overtime record", map[string]interface{}{
		"overtime_id": id,
	}, c)
	if err := o.OvertimeService.UpdateRecordOvertime(uint(id), &payload, c, tx); err != nil {
		helpers.MyLogger("error", "OvertimeManagement", "UpdateRecordOvertime", "controller", "error update overtime record", map[string]interface{}{
			"error": err.Error(),
		}, c)
		return helpers.ResponseErrorInternal(c, err)
	}
	return nil
}

// DeleteRecordOvertime deletes an overtime record
func (o *OvertimeController) DeleteRecordOvertime(c *fiber.Ctx) error {
	helpers.MyLogger("debug", "OvertimeManagement", "DeleteRecordOvertime", "controller", "start delete overtime record", nil, c)

	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		helpers.MyLogger("error", "OvertimeManagement", "DeleteRecordOvertime", "controller", "error parse overtime ID", map[string]interface{}{
			"error": err.Error(),
		}, c)
		return helpers.ResponseErrorBadRequest(c, "Invalid overtime ID", nil)
	}

	tx := database.ClientPostgres.Begin()
	defer tx.Rollback()

	helpers.MyLogger("debug", "OvertimeManagement", "DeleteRecordOvertime", "controller", "start calling service for delete overtime record", map[string]interface{}{
		"overtime_id": id,
	}, c)
	if err := o.OvertimeService.DeleteRecordOvertime(uint(id), c, tx); err != nil {
		helpers.MyLogger("error", "OvertimeManagement", "DeleteRecordOvertime", "controller", "error delete overtime record", map[string]interface{}{
			"error": err.Error(),
		}, c)
		return helpers.ResponseErrorInternal(c, err)
	}
	return nil
}
