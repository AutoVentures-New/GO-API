package configuration

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hubjob/api/app/adapters/area"
	"github.com/hubjob/api/handler/responses"
	"github.com/hubjob/api/model"
)

func GetCulturalFit(fiberCtx *fiber.Ctx) error {
	return responses.Success(fiberCtx, model.CulturalFit)
}

func ListAreas(fiberCtx *fiber.Ctx) error {
	areas, err := area.ListAreas(fiberCtx.UserContext())
	if err != nil {
		return responses.InternalServerError(fiberCtx, err)
	}

	return responses.Success(fiberCtx, areas)
}
