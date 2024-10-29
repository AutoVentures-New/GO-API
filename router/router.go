package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/trabalhe-conosco/api/handler"
	"github.com/trabalhe-conosco/api/middleware"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/")
	api.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.JSON(fiber.Map{"status": "success", "message": "Hello i'm ok!", "data": nil})
	})

	setupB2bRoute(api.Group("/b2b"))
}

func setupB2bRoute(router fiber.Router) {
	auth := router.Group("/auth")
	auth.Post("/register", handler.Register)
	auth.Post("/login", handler.Login)

	auth.Get("/me", middleware.Protected(), handler.Me)
}
