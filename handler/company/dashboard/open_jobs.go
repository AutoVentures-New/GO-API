package dashboard

import (
	"github.com/gofiber/fiber/v2"
	dashboard_adp "github.com/hubjob/api/app/adapters/company/dashboard"
	"github.com/hubjob/api/handler/responses"
	"github.com/hubjob/api/model"
)

func OpenJobs(fiberCtx *fiber.Ctx) error {
	user := fiberCtx.Locals("user").(model.User)

	jobs, err := dashboard_adp.OpenJobs(
		fiberCtx.UserContext(),
		user.CompanyID,
	)
	if err != nil {
		return responses.InternalServerError(fiberCtx, err)
	}

	return responses.Success(fiberCtx, jobs)
}

func Applications(fiberCtx *fiber.Ctx) error {
	user := fiberCtx.Locals("user").(model.User)

	applications, err := dashboard_adp.Applications(
		fiberCtx.UserContext(),
		user.CompanyID,
	)
	if err != nil {
		return responses.InternalServerError(fiberCtx, err)
	}

	return responses.Success(fiberCtx, applications)
}
