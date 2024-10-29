package middleware

import "github.com/gofiber/fiber/v2"

func ErrorHandler() fiber.ErrorHandler {
	return func(fiberCtx *fiber.Ctx, err error) error {
		return fiberCtx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal Server Error",
		})
	}
}
