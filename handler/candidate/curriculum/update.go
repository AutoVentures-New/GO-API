package curriculum

import (
	"github.com/gofiber/fiber/v2"
	curriculum_adp "github.com/hubjob/api/app/adapters/candidate/curriculum"
	"github.com/hubjob/api/handler/responses"
	"github.com/hubjob/api/model"
)

func UpdateCurriculum(fiberCtx *fiber.Ctx) error {
	candidate := fiberCtx.Locals("candidate").(model.Candidate)
	request := new(model.Curriculum)

	if err := fiberCtx.BodyParser(&request); err != nil {
		return responses.InvalidBodyRequest(fiberCtx, err)
	}

	request.CandidateID = candidate.ID

	curriculum, err := curriculum_adp.UpdateCurriculum(fiberCtx.UserContext(), *request)
	if err != nil {
		return responses.InternalServerError(fiberCtx, err)
	}

	return responses.Success(fiberCtx, curriculum)
}
