package public

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/hubjob/api/app/adapters/public"
	"github.com/hubjob/api/handler/responses"
)

func GetCompany(fiberCtx *fiber.Ctx) error {
	id := fiberCtx.Params("company_id")
	if len(id) == 0 {
		return responses.BadRequest(fiberCtx, "Params {company_id} is required")
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		return responses.BadRequest(fiberCtx, "Invalid params {company_id}")
	}

	company, err := public.GetCompany(fiberCtx.UserContext(), int64(idInt))
	if errors.Is(err, public.ErrNotFound) {
		return responses.NotFound(fiberCtx, err.Error())
	}

	if err != nil {
		return responses.InternalServerError(fiberCtx, err)
	}

	return responses.Success(fiberCtx, company)
}

func ListCompanies(fiberCtx *fiber.Ctx) error {
	filter := public.ListCompaniesFilter{}

	err := fiberCtx.QueryParser(&filter)
	if err != nil {
		return responses.BadRequest(fiberCtx, "Invalid filters")
	}

	companies, total, err := public.ListCompanies(
		fiberCtx.UserContext(),
		filter,
	)
	if err != nil {
		return responses.InternalServerError(fiberCtx, err)
	}

	return responses.Success(fiberCtx, map[string]any{
		"total":     total,
		"companies": companies,
	})
}
