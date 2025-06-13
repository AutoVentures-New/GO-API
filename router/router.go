package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/AutoVentures-New/GO-API/handler/responses"
	"github.com/AutoVentures-New/GO-API/middleware"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/")
	api.Get("/", func(fiberCtx *fiber.Ctx) error {
		return responses.Success(fiberCtx, "I'm OK")
	})

	protected := api.Group("/", middleware.Auth())
	setupContactDataRoute(protected.Group("/contact"))
}

func RouteNotFound() fiber.Handler {
	return func(fiberCtx *fiber.Ctx) error {
		return fiberCtx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "ROUTE_NOT_FOUND",
		})
	}
}
