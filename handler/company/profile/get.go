package profile

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hubjob/api/app/adapters/company/profile"
	"github.com/hubjob/api/handler/responses"
	"github.com/hubjob/api/model"
)

func GetCompany(fiberCtx *fiber.Ctx) error {
	user := fiberCtx.Locals("user").(model.User)

	company, err := profile.GetCompany(fiberCtx.UserContext(), user.CompanyID)
	if err != nil {
		return responses.InternalServerError(fiberCtx, err)
	}

	return responses.Success(fiberCtx, company)
}
