package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hubjob/api/handler/company"
)

func setupCompanyRoute(router fiber.Router) {
	auth := router.Group("/auth/company")

	auth.Post("/validate-cnpj-cpf", company.ValidateCnpjCpf)
	auth.Post("/send-email-validation", company.SendEmailValidation)
	auth.Post("/email-validation-code", company.ValidateEmailValidationCode)

	auth.Post("/create-account", company.CreateAccount)
	//auth.Post("/login", handler.Login)

	//auth.Get("/me", middleware.Protected(), handler.Me)
}
