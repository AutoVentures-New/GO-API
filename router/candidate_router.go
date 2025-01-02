package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hubjob/api/handler/candidate"
	"github.com/hubjob/api/middleware"
)

func setupCandidateRoute(router fiber.Router) {
	auth := router.Group("/auth")

	auth.Post("/validate-cpf", candidate.ValidateCpf)
	auth.Post("/send-email-validation", candidate.SendEmailValidation)
	auth.Post("/email-validation-code", candidate.ValidateEmailValidationCode)

	auth.Post("/create-account", candidate.CreateAccount)
	auth.Post("/login", candidate.Login)

	auth.Get("/me", middleware.Protected(), candidate.Me)
}
