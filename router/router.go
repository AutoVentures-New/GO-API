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

	//setupB2bRoute(api.Group("/b2b"))
	setupCompanyRoute(api.Group("/company"))
}

//func setupB2bRoute(router fiber.Router) {
//	auth := router.Group("/auth")
//	auth.Post("/register", handler.Register)
//	auth.Post("/login", handler.Login)
//
//	auth.Get("/me", middleware.Protected(), handler.Me)
//}
