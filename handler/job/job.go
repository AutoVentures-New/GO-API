package job

import (
	"github.com/gofiber/fiber/v2"
	job_adp "github.com/hubjob/api/app/adapters/job"
	"github.com/hubjob/api/handler/responses"
)

func ListJobs(fiberCtx *fiber.Ctx) error {
	filter := job_adp.Filter{}

	err := fiberCtx.QueryParser(&filter)
	if err != nil {
		return responses.BadRequest(fiberCtx, "Invalid filters")
	}

	if filter.Page <= 0 {
		filter.Page = 1
	}

	if filter.Size < 5 {
		filter.Size = 5
	}

	jobs, total, err := job_adp.ListJobs(
		fiberCtx.UserContext(),
		filter,
	)
	if err != nil {
		return responses.InternalServerError(fiberCtx, err)
	}

	nextPage := filter.Page
	if total > 0 && (total/filter.Size) >= filter.Page {
		nextPage = filter.Page + 1
	}

	previousPage := filter.Page - 1
	if previousPage <= 0 {
		previousPage = 1
	}

	return responses.Success(fiberCtx, map[string]any{
		"total":         total,
		"current_page":  filter.Page,
		"next_page":     nextPage,
		"previous_page": previousPage,
		"jobs":          jobs,
	})
}
