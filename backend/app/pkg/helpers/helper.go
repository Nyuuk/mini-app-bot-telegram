package helpers

import (
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type Payload interface {
	CustomErrorsMessage(validator.ValidationErrors) []map[string]string
}

func ValidateBody(payload Payload, c *fiber.Ctx) error {
	// log.Debug("Validating request body ", payload)
	// Logger.Debug().Interface("payload", payload).Msg("Validating request body")
	validate := validator.New()

	// Parse the request body
	if err := c.BodyParser(&payload); err != nil {
		// log.Error("Error parsing re quest body: ", err)
		// Logger.Error().Err(err).Msg("Error parsing request body")
		// ResponseErrorBadRequest(c, "Invalid payload", err)
		return err
	}
	// log.Debug("Validating after body parser ", payload)
	// Logger.Debug().Interface("payload", payload).Msg("Validating after body parser")

	// Validate the user struct
	if err := validate.Struct(payload); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		data := payload.CustomErrorsMessage(validationErrors)
		// log.Debug("Validation errors ", data)
		// Logger.Debug().Interface("validation_errors", data).Msg("Validation errors")
		// log.Debug("Validation errors ", validationErrors)
		return ErrorClient("Invalid payload ", fiber.StatusBadRequest, data)
		// return ErrorClient("Invalid payload ", fiber.StatusBadRequest, validationErrors)
	}

	return nil
}

func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
