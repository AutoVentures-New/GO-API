package job

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	job "github.com/hubjob/api/app/adapters/company/job/application"
	"github.com/hubjob/api/handler/responses"
)

type UpdateQuestionScoreRequest struct {
	QuestionID int64 `json:"question_id"`
	Score      int64 `json:"score"`
}

func UpdateQuestionScore(fiberCtx *fiber.Ctx) error {
	request := new(UpdateQuestionScoreRequest)

	if err := fiberCtx.BodyParser(&request); err != nil {
		return responses.InvalidBodyRequest(fiberCtx, err)
	}

	if request.Score > 5 {
		request.Score = 5
	}

	if request.QuestionID == 0 {
		return responses.BadRequest(fiberCtx, "Params {question_id} is required")
	}

	id := fiberCtx.Params("application_id")
	if len(id) == 0 {
		return responses.BadRequest(fiberCtx, "Params {application_id} is required")
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		return responses.BadRequest(fiberCtx, "Invalid params {application_id}")
	}

	err = job.UpdateQuestionScore(fiberCtx.UserContext(), int64(idInt), request.QuestionID, request.Score)
	if err != nil {
		return responses.InternalServerError(fiberCtx, err)
	}

	return responses.Success(fiberCtx, nil)
}
