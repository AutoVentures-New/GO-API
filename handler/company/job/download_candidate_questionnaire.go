package job

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	job "github.com/hubjob/api/app/adapters/company/job/application"
	"github.com/hubjob/api/handler/responses"
	"github.com/hubjob/api/model"
)

func DownloadCandidateQuestionnaire(fiberCtx *fiber.Ctx) error {
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

	file, err := job.DownloadCandidateQuestionnaire(fiberCtx.UserContext(), int64(jobIdInt), user.CompanyID, int64(idInt))
	if err != nil && !errors.Is(err, job.ErrApplicationNotFound) && !errors.Is(err, job.ErrFileNotFound) {
		return responses.InternalServerError(fiberCtx, err)
	}

	if errors.Is(err, job.ErrFileNotFound) || errors.Is(err, job.ErrApplicationNotFound) || file == nil {
		return responses.NotFound(fiberCtx, err.Error())
	}

	return responses.Download(fiberCtx, "file", file)
}
