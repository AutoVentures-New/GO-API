package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hubjob/api/handler/candidate"
	"github.com/hubjob/api/handler/candidate/curriculum"
	"github.com/hubjob/api/handler/candidate/profile"
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

	auth.Get("/me", middleware.ProtectedCandidate(false), candidate.Me)

	prof := router.Group("/profile", middleware.ProtectedCandidate(false))

	prof.Patch("", profile.UpdateCandidate)
	prof.Patch("update-email", profile.UpdateCandidateEmail)
	prof.Get("photo", profile.DownloadCandidatePhoto)
	prof.Post("photo", profile.UpdateCandidatePhoto)

	candidateCurriculum := router.Group("/curriculum", middleware.ProtectedCandidate(false))

	candidateCurriculum.Get("", curriculum.GetCurriculum)
	candidateCurriculum.Patch("", curriculum.UpdateCurriculum)

	job := router.Group("/job/application", middleware.ProtectedCandidate(false))

	job.Get("", candidate.ListApplications)
	job.Get("/:job_id", candidate.GetApplication)
	job.Delete("/:job_id", candidate.CanceledApplication)
	job.Post("/:job_id/start", candidate.StartApplication)

	stepsRoute := job.Group("/:job_id/steps", middleware.ProtectedCandidate(false))

	stepsRoute.Post("requirements", steps.SaveRequirements)
	stepsRoute.Post("job-questions", steps.SaveJobQuestions)
	stepsRoute.Post("cultural-fit", steps.SaveCulturalFit)
	stepsRoute.Post("questionnaire", steps.SaveQuestionnaire)
	stepsRoute.Post("candidate-video", steps.SaveCandidateVideo)
}
