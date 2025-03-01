package configuration

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hubjob/api/app/adapters/area"
	company_job_adp "github.com/hubjob/api/app/adapters/company/job"
	"github.com/hubjob/api/handler/responses"
	"github.com/hubjob/api/model"
	"strconv"
)

func GetCulturalFit(fiberCtx *fiber.Ctx) error {
	return responses.Success(fiberCtx, model.CulturalFit)
}

func GetQuestionnaireBehavioral(fiberCtx *fiber.Ctx) error {
	return responses.Success(fiberCtx, model.QuestionnaireBehavioral)
}

func GetQuestionnaireProfessional(fiberCtx *fiber.Ctx) error {
	return responses.Success(fiberCtx, model.QuestionnaireProfessional)
}

func ListAreas(fiberCtx *fiber.Ctx) error {
	areas, err := area.ListAreas(fiberCtx.UserContext())
	if err != nil {
		return responses.InternalServerError(fiberCtx, err)
	}

	return responses.Success(fiberCtx, areas)
}

func ListStateCitiesByCompany(fiberCtx *fiber.Ctx) error {
	id := fiberCtx.Params("company_id")
	if len(id) == 0 {
		return responses.BadRequest(fiberCtx, "Params {company_id} is required")
	}

	idInt, err := strconv.Atoi(id)
	if len(id) == 0 {
		return responses.BadRequest(fiberCtx, "Invalid params {company_id}")
	}

	jobs, err := company_job_adp.ListStateCities(
		fiberCtx.UserContext(),
		int64(idInt),
	)
	if err != nil {
		return responses.InternalServerError(fiberCtx, err)
	}

	return responses.Success(fiberCtx, jobs)
}
