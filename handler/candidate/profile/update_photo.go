package profile

import (
	"mime"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/hubjob/api/app/adapters/candidate/profile"
	"github.com/hubjob/api/handler/responses"
	"github.com/hubjob/api/model"
)

func UpdateCandidatePhoto(fiberCtx *fiber.Ctx) error {
	candidate := fiberCtx.Locals("candidate").(model.Candidate)

	photo, err := fiberCtx.FormFile("photo")
	if err != nil {
		return responses.InvalidBodyRequest(fiberCtx, err)
	}

	contentType := photo.Header.Get("Content-Type")
	if contentType == "" {
		contentType = mime.TypeByExtension(filepath.Ext(photo.Filename))
	}

	err = profile.UpdateCandidatePhoto(fiberCtx.UserContext(), candidate.ID, photo, contentType)
	if err != nil {
		return responses.InternalServerError(fiberCtx, err)
	}

	return responses.Success(fiberCtx, nil)
}
