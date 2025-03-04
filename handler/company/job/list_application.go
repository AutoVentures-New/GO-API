package job

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	job "github.com/hubjob/api/app/adapters/company/job/application"
	"github.com/hubjob/api/handler/responses"
	"github.com/hubjob/api/model"
)

func ListJobApplications(fiberCtx *fiber.Ctx) error {
	user := fiberCtx.Locals("user").(model.User)

	id := fiberCtx.Params("id")
	if len(id) == 0 {
		return responses.BadRequest(fiberCtx, "Params {id} is required")
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		return responses.BadRequest(fiberCtx, "Invalid params {id}")
	}

	request := new(job.ListJobApplicationsRequest)

	if err := fiberCtx.BodyParser(&request); err != nil {
		return responses.InvalidBodyRequest(fiberCtx, err)
	}

	request.FilterCandidateName = strings.TrimSpace(request.FilterCandidateName)

	applications, err := job.ListJobApplications(fiberCtx.UserContext(), user.CompanyID, int64(idInt), *request)
	if err != nil {
		return responses.InternalServerError(fiberCtx, err)
	}

	return responses.Success(fiberCtx, applications)
}
