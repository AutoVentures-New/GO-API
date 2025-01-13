package company

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	company_benefit_adp "github.com/hubjob/api/app/adapters/company/benefit"
	"github.com/hubjob/api/handler/responses"
	"github.com/hubjob/api/model"
)

type CreateBenefitRequest struct {
	Name string `json:"name"`
}

func CreateBenefit(fiberCtx *fiber.Ctx) error {
	user := fiberCtx.Locals("user").(model.User)
	request := new(CreateBenefitRequest)

	if err := fiberCtx.BodyParser(&request); err != nil {
		return responses.InvalidBodyRequest(fiberCtx, err)
	}

	benefit, err := company_benefit_adp.CreateBenefit(
		fiberCtx.UserContext(),
		model.Benefit{
			Name:      request.Name,
			CompanyID: user.CompanyID,
		},
	)
	if errors.Is(err, company_benefit_adp.ErrBenefitAlreadyExists) {
		return responses.Conflict(fiberCtx, "COMPANY|BENEFIT|ALREADY_EXISTS", "Benefit already exists")
	}

	if err != nil {
		return responses.InternalServerError(fiberCtx, err)
	}

	return responses.Success(fiberCtx, benefit)
}

func ListBenefits(fiberCtx *fiber.Ctx) error {
	user := fiberCtx.Locals("user").(model.User)

	benefits, err := company_benefit_adp.ListBenefits(
		fiberCtx.UserContext(),
		user.CompanyID,
	)
	if err != nil {
		return responses.InternalServerError(fiberCtx, err)
	}

	return responses.Success(fiberCtx, benefits)
}

func GetBenefit(fiberCtx *fiber.Ctx) error {
	user := fiberCtx.Locals("user").(model.User)

	id := fiberCtx.Params("id")
	if len(id) == 0 {
		return responses.BadRequest(fiberCtx, "Params {id} is required")
	}

	idInt, err := strconv.Atoi(id)
	if len(id) == 0 {
		return responses.BadRequest(fiberCtx, "Invalid params {id}")
	}

	benefit, err := company_benefit_adp.GetBenefit(
		fiberCtx.UserContext(),
		int64(idInt),
		user.CompanyID,
	)
	if errors.Is(err, company_benefit_adp.ErrNotFound) {
		return responses.NotFound(fiberCtx, err.Error())
	}

	if err != nil {
		return responses.InternalServerError(fiberCtx, err)
	}

	return responses.Success(fiberCtx, benefit)
}

type UpdateBenefitRequest struct {
	Name string `json:"name"`
}

func UpdateBenefit(fiberCtx *fiber.Ctx) error {
	user := fiberCtx.Locals("user").(model.User)
	request := new(UpdateBenefitRequest)

	if err := fiberCtx.BodyParser(&request); err != nil {
		return responses.InvalidBodyRequest(fiberCtx, err)
	}

	id := fiberCtx.Params("id")
	if len(id) == 0 {
		return responses.BadRequest(fiberCtx, "Params {id} is required")
	}

	idInt, err := strconv.Atoi(id)
	if len(id) == 0 {
		return responses.BadRequest(fiberCtx, "Invalid params {id}")
	}

	benefit, err := company_benefit_adp.GetBenefit(
		fiberCtx.UserContext(),
		int64(idInt),
		user.CompanyID,
	)
	if errors.Is(err, company_benefit_adp.ErrNotFound) {
		return responses.NotFound(fiberCtx, err.Error())
	}

	if err != nil {
		return responses.InternalServerError(fiberCtx, err)
	}

	if request.Name == benefit.Name {
		return responses.Success(fiberCtx, benefit)
	}

	benefit.Name = request.Name
	benefit, err = company_benefit_adp.UpdateBenefit(
		fiberCtx.UserContext(),
		benefit,
	)
	if err != nil {
		return responses.InternalServerError(fiberCtx, err)
	}

	return responses.Success(fiberCtx, benefit)
}

func DeleteBenefit(fiberCtx *fiber.Ctx) error {
	user := fiberCtx.Locals("user").(model.User)

	id := fiberCtx.Params("id")
	if len(id) == 0 {
		return responses.BadRequest(fiberCtx, "Params {id} is required")
	}

	idInt, err := strconv.Atoi(id)
	if len(id) == 0 {
		return responses.BadRequest(fiberCtx, "Invalid params {id}")
	}

	err = company_benefit_adp.DeleteBenefit(
		fiberCtx.UserContext(),
		int64(idInt),
		user.CompanyID,
	)
	if err != nil {
		return responses.InternalServerError(fiberCtx, err)
	}

	return responses.Success(fiberCtx, nil)
}
