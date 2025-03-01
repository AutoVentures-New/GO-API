package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hubjob/api/handler/company"
	"github.com/hubjob/api/handler/company/profile"
	"github.com/hubjob/api/middleware"
)

func setupCompanyRoute(router fiber.Router) {
	auth := router.Group("/auth")

	auth.Post("/validate-cnpj-cpf", company.ValidateCnpjCpf)
	auth.Post("/send-email-validation", company.SendEmailValidation)
	auth.Post("/email-validation-code", company.ValidateEmailValidationCode)

	auth.Post("/create-account", company.CreateAccount)
	auth.Post("/login", company.Login)

	auth.Get("/me", middleware.ProtectedCompany(), company.Me)

	benefit := router.Group("/benefit", middleware.ProtectedCompany())

	benefit.Post("", company.CreateBenefit)
	benefit.Get("", company.ListBenefits)
	benefit.Get("/:id", company.GetBenefit)
	benefit.Patch("/:id", company.UpdateBenefit)
	benefit.Delete("/:id", company.DeleteBenefit)

	questionnaire := router.Group("/questionnaire", middleware.ProtectedCompany())

	questionnaire.Post("", company.CreateQuestionnaire)
	questionnaire.Get("", company.ListQuestionnaires)
	questionnaire.Get("/:id", company.GetQuestionnaire)
	questionnaire.Patch("/:id", company.UpdateQuestionnaire)
	questionnaire.Delete("/:id", company.DeleteQuestionnaire)

	question := questionnaire.Group("/:questionnaire_id/question", middleware.ProtectedCompany())

	question.Post("", company.CreateQuestion)
	question.Get("", company.ListQuestions)
	question.Patch("/:id", company.UpdateQuestion)
	question.Delete("/:id", company.DeleteQuestion)

	job := router.Group("/job", middleware.ProtectedCompany())

	job.Post("", company.CreateJob)
	job.Get("", company.ListJobs)
	job.Get("/:id", company.GetJob)
	job.Patch("/:id", company.UpdateJob)
	job.Delete("/:id", company.DeleteJob)

	router.Get("/cultural-fit", middleware.ProtectedCompany(), company.GetLastCulturalFit)

	prof := router.Group("/profile", middleware.ProtectedCompany())

	prof.Get("", profile.GetCompany)
	prof.Patch("", profile.UpdateCompany)
}
