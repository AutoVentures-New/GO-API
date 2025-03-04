package curriculum

import (
	"github.com/gofiber/fiber/v2"
	curriculum_adp "github.com/hubjob/api/app/adapters/candidate/curriculum"
	"github.com/hubjob/api/handler/responses"
	"github.com/hubjob/api/model"
)

func GetCurriculum(fiberCtx *fiber.Ctx) error {
	candidate := fiberCtx.Locals("candidate").(model.Candidate)

	curriculum, err := curriculum_adp.GetCurriculum(fiberCtx.UserContext(), candidate.ID)
	if err != nil {
		return responses.InternalServerError(fiberCtx, err)
	}

	return responses.Success(fiberCtx, curriculum)
}
