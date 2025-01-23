package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hubjob/api/handler/job"
)

func setupJobRoute(router fiber.Router) {
	router.Get("", job.ListJobs)
}
