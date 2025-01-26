package steps

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	candidate_job_adp "github.com/hubjob/api/app/adapters/candidate/job"
	"github.com/hubjob/api/app/adapters/candidate/job/steps"
	"github.com/hubjob/api/handler/responses"
	"github.com/hubjob/api/model"
)

type SaveJobQuestionsRequest struct {
	JobID     int64                       `json:"job_id"`
	Questions []model.ApplicationQuestion `json:"questions"`
}

func SaveJobQuestions(fiberCtx *fiber.Ctx) error {
	candidate := fiberCtx.Locals("candidate").(model.Candidate)
	request := new(SaveJobQuestionsRequest)

	if err := fiberCtx.BodyParser(&request); err != nil {
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

	if application.Status != model.FILLING || application.CurrentStep != model.JOB_QUESTIONS {
		return responses.BadRequest(fiberCtx, "Invalid step")
	}

	application.Questions = request.Questions

	application, err = steps.SaveJobQuestions(fiberCtx.UserContext(), application)
	if err != nil {
		return responses.InternalServerError(fiberCtx, err)
	}

	return responses.Success(fiberCtx, application)
}
