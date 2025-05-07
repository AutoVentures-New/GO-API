package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hubjob/api/handler/company"
	job2 "github.com/hubjob/api/handler/company/job"
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
	auth.Post("/create-user-password", company.CreateUserPassword)

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

	question := questionnaire.Group("/:questionnaire_id/question")

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

	applications := job.Group("/:id/application")

	applications.Post("", job2.ListJobApplications)

	application := applications.Group("/:application_id")

	application.Get("", job2.GetJobApplication)
	application.Patch("/candidate-video-score", job2.UpdateCandidateVideoScore)
	application.Patch("/question-score", job2.UpdateQuestionScore)
	application.Patch("/evaluate-candidate", job2.EvaluateCandidate)
	application.Get("/candidate-video", job2.DownloadCandidateVideo)
	application.Get("/candidate-photo", job2.DownloadCandidatePhoto)
	application.Get("/candidate-questionnaire-result", job2.DownloadCandidateQuestionnaire)

	router.Get("/cultural-fit", middleware.ProtectedCompany(), company.GetLastCulturalFit)

	prof := router.Group("/profile", middleware.ProtectedCompany())

	prof.Get("", profile.GetCompany)
	prof.Patch("", profile.UpdateCompany)

	users := router.Group("/users", middleware.ProtectedCompany())

	users.Post("", company.CreateUser)
	users.Get("", company.ListUsers)
	users.Patch("/:id", company.UpdateUser)
	users.Delete("/:id", company.DeleteUser)
}
