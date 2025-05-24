package public

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	public_adp "github.com/hubjob/api/app/adapters/public"
	"github.com/hubjob/api/handler/responses"
	"github.com/hubjob/api/model"
	"strconv"
)

func ListJobs(fiberCtx *fiber.Ctx) error {
	filter := public_adp.Filter{}

	err := fiberCtx.QueryParser(&filter)
	if err != nil {
		return responses.BadRequest(fiberCtx, "Invalid filters")
	}

	if filter.Page <= 0 {
		filter.Page = 1
	}

	if filter.Size < 20 {
		filter.Size = 20
	}

	jobs, total, err := public_adp.ListJobs(
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

func GetJob(fiberCtx *fiber.Ctx) error {
	id := fiberCtx.Params("job_id")
	if len(id) == 0 {
		return responses.BadRequest(fiberCtx, "Params {job_id} is required")
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		return responses.BadRequest(fiberCtx, "Invalid params {job_id}")
	}

	candidate := model.Candidate{}

	candidateLocal := fiberCtx.Locals("candidate")
	if candidateLocal != nil && candidateLocal != "" {
		candidate = candidateLocal.(model.Candidate)
	}

	job, err := public_adp.GetJob(fiberCtx.UserContext(), int64(idInt), candidate.ID)
	if errors.Is(err, public_adp.ErrJobNotFound) {
		return responses.NotFound(fiberCtx, err.Error())
	}

	if err != nil {
		return responses.InternalServerError(fiberCtx, err)
	}

	return responses.Success(fiberCtx, job)
}
