package profile

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/hubjob/api/app/adapters/candidate/profile"
	"github.com/hubjob/api/handler/responses"
	"github.com/hubjob/api/model"
)

func DownloadCandidatePhoto(fiberCtx *fiber.Ctx) error {
	candidate := fiberCtx.Locals("candidate").(model.Candidate)

	photo, err := profile.DownloadCandidatePhoto(fiberCtx.UserContext(), candidate.ID)
	if err != nil && !errors.Is(err, profile.ErrPhotoNotFound) {
		return responses.InternalServerError(fiberCtx, err)
	}

	if errors.Is(err, profile.ErrPhotoNotFound) || photo == nil {
		return responses.NotFound(fiberCtx, "photo not found")
	}

	return responses.Download(fiberCtx, "photo", photo)
}
