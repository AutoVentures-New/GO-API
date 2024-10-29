package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func Me(fiberCtx *fiber.Ctx) error {
	token := fiberCtx.Locals("user").(*jwt.Token)

	return fiberCtx.JSON(fiber.Map{"data": token.Claims})
}
