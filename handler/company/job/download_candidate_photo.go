package job

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	job "github.com/hubjob/api/app/adapters/company/job/application"
	"github.com/hubjob/api/handler/responses"
	"github.com/hubjob/api/model"
)

func DownloadCandidatePhoto(fiberCtx *fiber.Ctx) error {
	user := fiberCtx.Locals("user").(model.User)

	id := fiberCtx.Params("application_id")
	if len(id) == 0 {
		return responses.BadRequest(fiberCtx, "Params {application_id} is required")
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		return responses.BadRequest(fiberCtx, "Invalid params {application_id}")
	}

	photo, err := job.DownloadCandidatePhoto(fiberCtx.UserContext(), user.CompanyID, int64(idInt))
	if err != nil && !errors.Is(err, job.ErrPhotoNotFound) {
		return responses.InternalServerError(fiberCtx, err)
	}

	if errors.Is(err, job.ErrPhotoNotFound) || photo == nil {
		return responses.NotFound(fiberCtx, "photo not found")
	}

	return responses.Download(fiberCtx, "photo", photo)
}
