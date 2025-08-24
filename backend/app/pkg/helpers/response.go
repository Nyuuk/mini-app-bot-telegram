package helpers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func Response(c *fiber.Ctx, status int, message string, data interface{}) error {
	return c.Status(status).JSON(fiber.Map{
		"code":    status,
		"message": message,
		"data":    data,
	})
}

func ResponseErrorInternal(c *fiber.Ctx, err any) error {

	return Response(c, fiber.StatusInternalServerError, "Internal server error", nil)
}

func ResponseErrorBadRequest(c *fiber.Ctx, message string, err any) error {
	log.Error("Bad request: ", err)
	return Response(c, fiber.StatusBadRequest, message, err)
}

func ResponseErrorNotFound(c *fiber.Ctx, err any) error {
	log.Error("Not found: ", err)
	return Response(c, fiber.StatusNotFound, "Not found", nil)
}
