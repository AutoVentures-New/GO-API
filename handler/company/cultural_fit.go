package company

import (
	"github.com/gofiber/fiber/v2"
	company_job_adp "github.com/hubjob/api/app/adapters/company/job"
	"github.com/hubjob/api/handler/responses"
	"github.com/hubjob/api/model"
)

func GetLastCulturalFit(fiberCtx *fiber.Ctx) error {
	user := fiberCtx.Locals("user").(model.User)

	job, err := company_job_adp.GetLastCulturalFit(
		fiberCtx.UserContext(),
		user.CompanyID,
	)
	if err != nil {
		return responses.InternalServerError(fiberCtx, err)
	}

	return responses.Success(fiberCtx, job)
}
