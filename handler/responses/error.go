package responses

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func BadRequest(
	fiberCtx *fiber.Ctx,
	error string,
) error {
	return fiberCtx.Status(fiber.StatusBadRequest).
		JSON(fiber.Map{
			"message": "Invalid Request",
			"error":   error,
		})
}

func InvalidBodyRequest(
	fiberCtx *fiber.Ctx,
	err error,
) error {
	return fiberCtx.Status(fiber.StatusUnprocessableEntity).
		JSON(fiber.Map{
			"message": "Invalid Body Request",
			"error":   err.Error(),
		})
}

func NotFound(
	fiberCtx *fiber.Ctx,
	message string,
) error {
	return fiberCtx.Status(fiber.StatusNotFound).
		JSON(fiber.Map{
			"message": message,
			"error":   "NOT_FOUND",
		})
}

func MethodNotAllowed(
	fiberCtx *fiber.Ctx,
) error {
	return fiberCtx.Status(fiber.StatusMethodNotAllowed).
		JSON(fiber.Map{
			"message": "Method Not Allowed",
		})
}

func Unauthorized(
	fiberCtx *fiber.Ctx,
) error {
	return fiberCtx.Status(fiber.StatusUnauthorized).
		JSON(fiber.Map{
			"message": "Invalid email and password",
		})
}

func Forbidden(
	fiberCtx *fiber.Ctx,
) error {
	return fiberCtx.Status(fiber.StatusForbidden).
		JSON(fiber.Map{
			"message": "Invalid or expired token",
		})
}

func Conflict(
	fiberCtx *fiber.Ctx,
	error string,
	message string,
) error {
	return fiberCtx.Status(fiber.StatusConflict).
		JSON(fiber.Map{
			"message": message,
			"error":   error,
		})
}

func InternalServerError(
	fiberCtx *fiber.Ctx,
	err error,
) error {
	logrus.WithError(err).Error("Internal Server Error")

	return fiberCtx.Status(fiber.StatusInternalServerError).
		JSON(fiber.Map{
			"message": "Internal Server Error",
		})
}
