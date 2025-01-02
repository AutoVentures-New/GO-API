package middleware

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/hubjob/api/handler/responses"
	"github.com/sirupsen/logrus"
)

func ErrorHandler() fiber.ErrorHandler {
	return func(fiberCtx *fiber.Ctx, err error) error {
		var fiberErr *fiber.Error

		ok := errors.Is(err, fiberErr)
		if !ok {
			return responses.InternalServerError(fiberCtx, err)
		}

		if fiberErr.Code == fiber.StatusMethodNotAllowed {
			return responses.MethodNotAllowed(fiberCtx)
		}

		if fiberErr.Code == fiber.StatusNotFound {
			return responses.NotFound(fiberCtx, "Route Not Found")
		}

		logrus.WithError(err).
			WithFields(map[string]interface{}{
				"fiber_error_code":    fiberErr.Code,
				"fiber_error_message": fiberErr.Message,
			}).Error("Error not map")

		return responses.InternalServerError(fiberCtx, err)
	}
}
