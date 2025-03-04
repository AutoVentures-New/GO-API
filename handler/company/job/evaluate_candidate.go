package job

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	job "github.com/hubjob/api/app/adapters/company/job/application"
	"github.com/hubjob/api/handler/responses"
	"github.com/hubjob/api/model"
)

type EvaluateCandidateRequest struct {
	Status string `json:"status"`
}

func EvaluateCandidate(fiberCtx *fiber.Ctx) error {
	user := fiberCtx.Locals("user").(model.User)

	request := new(EvaluateCandidateRequest)

	if err := fiberCtx.BodyParser(&request); err != nil {
		return responses.InvalidBodyRequest(fiberCtx, err)
	}

	if request.Status != model.APPROVED && request.Status != model.REPROVED && request.Status != model.WAITING_EVALUATION {
		return responses.InvalidBodyRequest(fiberCtx, errors.New("invalid status"))
	}

	id := fiberCtx.Params("application_id")
	if len(id) == 0 {
		return responses.BadRequest(fiberCtx, "Params {application_id} is required")
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		return responses.BadRequest(fiberCtx, "Invalid params {application_id}")
	}

	err = job.EvaluateCandidate(fiberCtx.UserContext(), int64(idInt), user.CompanyID, request.Status)
	if err != nil {
		return responses.InternalServerError(fiberCtx, err)
	}

	return responses.Success(fiberCtx, nil)
}
