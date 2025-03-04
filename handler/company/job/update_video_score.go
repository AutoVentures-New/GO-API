package job

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	job "github.com/hubjob/api/app/adapters/company/job/application"
	"github.com/hubjob/api/handler/responses"
)

type UpdateCandidateVideoScoreRequest struct {
	Score int64 `json:"score"`
}

func UpdateCandidateVideoScore(fiberCtx *fiber.Ctx) error {
	request := new(UpdateCandidateVideoScoreRequest)

	if err := fiberCtx.BodyParser(&request); err != nil {
		return responses.InvalidBodyRequest(fiberCtx, err)
	}

	id := fiberCtx.Params("application_id")
	if len(id) == 0 {
		return responses.BadRequest(fiberCtx, "Params {application_id} is required")
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		return responses.BadRequest(fiberCtx, "Invalid params {application_id}")
	}

	err = job.UpdateCandidateVideoScore(fiberCtx.UserContext(), int64(idInt), request.Score)
	if err != nil {
		return responses.InternalServerError(fiberCtx, err)
	}

	return responses.Success(fiberCtx, nil)
}
