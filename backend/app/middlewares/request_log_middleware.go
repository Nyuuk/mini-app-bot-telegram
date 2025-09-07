package middlewares

import (
	"strconv"

	"github.com/Nyuuk/mini-app-bot-telegram/backend/app/entities"
	"github.com/Nyuuk/mini-app-bot-telegram/backend/app/pkg/helpers"
	"github.com/Nyuuk/mini-app-bot-telegram/backend/app/repositories"
	"github.com/gofiber/fiber/v2"
)

// RequestLogMiddleware logs all HTTP requests to the database
func RequestLogMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Next()

		var userID uint = 0
		if user := c.Locals("user_id"); user != nil {
			switch v := user.(type) {
			case uint:
				userID = v
			case int:
				userID = uint(v)
			case string:
				// convert string to uint
				if parsedID, err := strconv.ParseUint(v, 10, 32); err == nil {
					userID = uint(parsedID)
				}
			}
		}

		reqBody := ""
		if c.Method() != "GET" && c.Method() != "HEAD" {
			body := c.Body()
			if len(body) > 0 {
				reqBody = string(body)
			}
		}

		respBody := string(c.Response().Body())

		log := entities.LogRequest{
			UserID:       userID,
			Method:       c.Method(),
			Endpoint:     c.Path(),
			RequestBody:  reqBody,
			ResponseBody: respBody,
		}
		repo := repositories.LogRequestRepository{}
		s, err := repo.CreateLogRequest(log)
		if err != nil {
			helpers.MyLogger("error", "RequestLogMiddleware", "CreateLogRequest", "middleware", "error creating log request", map[string]interface{}{
				"error": err.Error(),
			}, c)
			return nil
		}
		c.Locals("log_request_id", s)
		return nil
	}
}
