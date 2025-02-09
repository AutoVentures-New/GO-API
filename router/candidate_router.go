package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hubjob/api/handler/candidate"
	"github.com/hubjob/api/handler/candidate/steps"
	"github.com/hubjob/api/middleware"
)

func setupCandidateRoute(router fiber.Router) {
	auth := router.Group("/auth")

	auth.Post("/validate-cpf", candidate.ValidateCpf)
	auth.Post("/send-email-validation", candidate.SendEmailValidation)
	auth.Post("/email-validation-code", candidate.ValidateEmailValidationCode)

	auth.Post("/create-account", candidate.CreateAccount)
	auth.Post("/login", candidate.Login)

	auth.Get("/me", middleware.ProtectedCandidate(), candidate.Me)

	job := router.Group("/job/application", middleware.ProtectedCandidate())

	job.Post("start", candidate.StartApplication)
	job.Get("/:job_id", candidate.GetApplication)
	job.Delete("/:job_id", candidate.CanceledApplication)

	stepsRoute := router.Group("/job/application/:job_id/steps", middleware.ProtectedCandidate())

	stepsRoute.Post("requirements", steps.SaveRequirements)
	stepsRoute.Post("job-questions", steps.SaveJobQuestions)
	stepsRoute.Post("cultural-fit", steps.SaveCulturalFit)
	stepsRoute.Post("questionnaire", steps.SaveQuestionnaire)
}
