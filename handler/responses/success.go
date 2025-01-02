package responses

import (
	"github.com/gofiber/fiber/v2"
)

func Success(
	fiberCtx *fiber.Ctx,
	data any,
) error {
	jsonData := make(map[string]string)

	if data == nil {
		data = jsonData
	}

	return fiberCtx.Status(fiber.StatusOK).
		JSON(fiber.Map{
			"data": data,
		})
}
