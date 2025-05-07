package job

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	job "github.com/hubjob/api/app/adapters/company/job/application"
	"github.com/hubjob/api/handler/responses"
	"github.com/hubjob/api/model"
)

func GetJobApplication(fiberCtx *fiber.Ctx) error {
	user := fiberCtx.Locals("user").(model.User)

	jobId := fiberCtx.Params("id")
	if len(jobId) == 0 {
		return responses.BadRequest(fiberCtx, "Params {job_id} is required")
	}

	jobIdInt, err := strconv.Atoi(jobId)
	if err != nil {
		return responses.BadRequest(fiberCtx, "Invalid params {job_id}")
	}

	id := fiberCtx.Params("application_id")
	if len(id) == 0 {
		return responses.BadRequest(fiberCtx, "Params {application_id} is required")
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		return responses.BadRequest(fiberCtx, "Invalid params {application_id}")
	}

	application, err := job.GetApplication(fiberCtx.UserContext(), user.CompanyID, int64(jobIdInt), int64(idInt))
	if err != nil && errors.Is(err, job.ErrApplicationNotFound) {
		return responses.NotFound(fiberCtx, err.Error())
	}

	if err != nil {
		return responses.InternalServerError(fiberCtx, err)
	}

	return responses.Success(fiberCtx, application)
}
