package job

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	job "github.com/hubjob/api/app/adapters/company/job/application"
	"github.com/hubjob/api/handler/responses"
)

func DownloadCandidateVideo(fiberCtx *fiber.Ctx) error {
	id := fiberCtx.Params("application_id")
	if len(id) == 0 {
		return responses.BadRequest(fiberCtx, "Params {application_id} is required")
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		return responses.BadRequest(fiberCtx, "Invalid params {application_id}")
	}

	video, err := job.DownloadCandidateVideo(fiberCtx.UserContext(), int64(idInt))
	if err != nil && !errors.Is(err, job.ErrVideoNotFound) {
		return responses.InternalServerError(fiberCtx, err)
	}

	if errors.Is(err, job.ErrVideoNotFound) || video == nil {
		return responses.NotFound(fiberCtx, "video not found")
	}

	return responses.Download(fiberCtx, "video", video)
}
