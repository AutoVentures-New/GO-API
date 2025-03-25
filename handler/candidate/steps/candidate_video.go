package steps

import (
	"errors"
	"mime"
	"path/filepath"
	"strconv"

	"github.com/gofiber/fiber/v2"
	candidate_job_adp "github.com/hubjob/api/app/adapters/candidate/job"
	"github.com/hubjob/api/app/adapters/candidate/job/steps"
	"github.com/hubjob/api/handler/responses"
	"github.com/hubjob/api/model"
)

func SaveCandidateVideo(fiberCtx *fiber.Ctx) error {
	candidate := fiberCtx.Locals("candidate").(model.Candidate)

	video, err := fiberCtx.FormFile("video")
	if err != nil {
		return responses.InvalidBodyRequest(fiberCtx, err)
	}

	id := fiberCtx.Params("job_id")
	if len(id) == 0 {
		return responses.BadRequest(fiberCtx, "Params {job_id} is required")
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		return responses.BadRequest(fiberCtx, "Invalid params {job_id}")
	}

	application, err := candidate_job_adp.GetJobApplication(
		fiberCtx.UserContext(),
		int64(idInt),
		candidate.ID,
	)
	if errors.Is(err, candidate_job_adp.ErrApplicationNotFound) {
		return responses.NotFound(fiberCtx, err.Error())
	}

	if err != nil {
		return responses.InternalServerError(fiberCtx, err)
	}

	if application.Status != model.FILLING || application.CurrentStep != model.CANDIDATE_VIDEO {
		return responses.BadRequest(fiberCtx, "Invalid step")
	}

	contentType := video.Header.Get("Content-Type")
	if contentType == "" {
		contentType = mime.TypeByExtension(filepath.Ext(video.Filename))
	}

	application, err = steps.SaveCandidateVideo(fiberCtx.UserContext(), application, video, contentType)
	if err != nil {
		return responses.InternalServerError(fiberCtx, err)
	}

	application, err = candidate_job_adp.GetJobApplication(
		fiberCtx.UserContext(),
		int64(idInt),
		candidate.ID,
	)
	if errors.Is(err, candidate_job_adp.ErrApplicationNotFound) {
		return responses.NotFound(fiberCtx, err.Error())
	}

	return responses.Success(fiberCtx, application)
}
