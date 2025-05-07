package profile

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hubjob/api/app/adapters/company/profile"
	"github.com/hubjob/api/handler/responses"
	"github.com/hubjob/api/model"
)

type UpdateCompanyRequest struct {
	Name        string  `json:"name"`
	CNPJ        string  `json:"cnpj"`
	Description *string `json:"description"`
}

func UpdateCompany(fiberCtx *fiber.Ctx) error {
	user := fiberCtx.Locals("user").(model.User)

	request := new(UpdateCompanyRequest)

	if err := fiberCtx.BodyParser(&request); err != nil {
		return responses.InvalidBodyRequest(fiberCtx, err)
	}

	company, err := profile.GetCompany(fiberCtx.UserContext(), user.CompanyID)
	if err != nil {
		return responses.InternalServerError(fiberCtx, err)
	}

	company.Name = request.Name
	company.CNPJ = request.CNPJ
	company.Description = request.Description

	company, err = profile.UpdateInfo(fiberCtx.UserContext(), company)
	if err != nil {
		return responses.InternalServerError(fiberCtx, err)
	}

	return responses.Success(fiberCtx, company)
}
