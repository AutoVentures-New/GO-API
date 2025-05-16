package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hubjob/api/handler/public"
)

func setupPublicRoute(router fiber.Router) {
	router.Get("list-jobs", public.ListJobs)
	router.Get("job-details/:job_id", public.GetJob)
	router.Get("company-details/:company_id", public.GetCompany)
	router.Get("list-companies", public.ListCompanies)
	router.Post("forgot-password", public.ForgotPassword)
	router.Post("change-password", public.ChangePassword)
}
