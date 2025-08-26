package helpers

import (
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type Payload interface {
	CustomErrorsMessage(validator.ValidationErrors) []map[string]string
}

func ValidateBody(payload Payload, c *fiber.Ctx) error {
	log.Debug("Validating request body ", payload)
	validate := validator.New()

	// Parse the request body
	if err := c.BodyParser(&payload); err != nil {
		log.Error("Error parsing request body: ", err)
		return Response(c, fiber.StatusBadRequest, "Invalid payload", nil)
	}
	log.Debug("Validating after body parser ", payload)

	// Validate the user struct
	if err := validate.Struct(payload); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		data := payload.CustomErrorsMessage(validationErrors)
		log.Debug("Validation errors ", data)
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
