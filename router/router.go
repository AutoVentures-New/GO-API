package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hubjob/api/handler/responses"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/")
	api.Get("/", func(fiberCtx *fiber.Ctx) error {
		return responses.Success(fiberCtx, "I'm OK")
	})

	setupCompanyRoute(api.Group("/company"))
	setupCandidateRoute(api.Group("/candidate"))
	setupConfigurationRoute(api.Group("/configuration"))
}
