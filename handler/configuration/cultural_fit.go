package configuration

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hubjob/api/handler/responses"
	"github.com/hubjob/api/model"
)

func GetCulturalFit(fiberCtx *fiber.Ctx) error {
	return responses.Success(fiberCtx, model.CulturalFit)
}
