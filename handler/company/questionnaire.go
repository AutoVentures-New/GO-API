package company

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	company_questionnaire_adp "github.com/hubjob/api/app/adapters/company/questionnaire"
	"github.com/hubjob/api/handler/responses"
	"github.com/hubjob/api/model"
)

type CreateQuestionnaireRequest struct {
	Name string `json:"name"`
}

func CreateQuestionnaire(fiberCtx *fiber.Ctx) error {
	user := fiberCtx.Locals("user").(model.User)
	request := new(CreateQuestionnaireRequest)

	if err := fiberCtx.BodyParser(&request); err != nil {
		return responses.InvalidBodyRequest(fiberCtx, err)
	}

	questionnaire, err := company_questionnaire_adp.CreateQuestionnaire(
		fiberCtx.UserContext(),
		model.Questionnaire{
			Name:      request.Name,
			CompanyID: user.CompanyID,
		},
	)
	if errors.Is(err, company_questionnaire_adp.ErrQuestionnaireAlreadyExists) {
		return responses.Conflict(fiberCtx, "COMPANY|QUESTIONNAIRE|ALREADY_EXISTS", "Questionnaire already exists")
	}

	if err != nil {
		return responses.InternalServerError(fiberCtx, err)
	}

	return responses.Success(fiberCtx, questionnaire)
}

func ListQuestionnaires(fiberCtx *fiber.Ctx) error {
	user := fiberCtx.Locals("user").(model.User)

	questionnaires, err := company_questionnaire_adp.ListQuestionnaires(
		fiberCtx.UserContext(),
		user.CompanyID,
	)
	if err != nil {
		return responses.InternalServerError(fiberCtx, err)
	}

	return responses.Success(fiberCtx, questionnaires)
}

func GetQuestionnaire(fiberCtx *fiber.Ctx) error {
	user := fiberCtx.Locals("user").(model.User)

	id := fiberCtx.Params("id")
	if len(id) == 0 {
		return responses.BadRequest(fiberCtx, "Params {id} is required")
	}

	idInt, err := strconv.Atoi(id)
	if len(id) == 0 {
		return responses.BadRequest(fiberCtx, "Invalid params {id}")
	}

	questionnaire, err := company_questionnaire_adp.GetQuestionnaire(
		fiberCtx.UserContext(),
		int64(idInt),
		user.CompanyID,
	)
	if errors.Is(err, company_questionnaire_adp.ErrQuestionnaireNotFound) {
		return responses.NotFound(fiberCtx, err.Error())
	}

	if err != nil {
		return responses.InternalServerError(fiberCtx, err)
	}

	return responses.Success(fiberCtx, questionnaire)
}

type UpdateQuestionnaireRequest struct {
	Name string `json:"name"`
}

func UpdateQuestionnaire(fiberCtx *fiber.Ctx) error {
	user := fiberCtx.Locals("user").(model.User)
	request := new(UpdateQuestionnaireRequest)

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

	questionnaire, err := company_questionnaire_adp.GetQuestionnaire(
		fiberCtx.UserContext(),
		int64(idInt),
		user.CompanyID,
	)
	if errors.Is(err, company_questionnaire_adp.ErrQuestionnaireNotFound) {
		return responses.NotFound(fiberCtx, err.Error())
	}

	if err != nil {
		return responses.InternalServerError(fiberCtx, err)
	}

	if request.Name == questionnaire.Name {
		return responses.Success(fiberCtx, questionnaire)
	}

	questionnaire.Name = request.Name
	questionnaire, err = company_questionnaire_adp.UpdateQuestionnaire(
		fiberCtx.UserContext(),
		questionnaire,
	)
	if err != nil {
		return responses.InternalServerError(fiberCtx, err)
	}

	return responses.Success(fiberCtx, questionnaire)
}

func DeleteQuestionnaire(fiberCtx *fiber.Ctx) error {
	user := fiberCtx.Locals("user").(model.User)

	id := fiberCtx.Params("id")
	if len(id) == 0 {
		return responses.BadRequest(fiberCtx, "Params {id} is required")
	}

	idInt, err := strconv.Atoi(id)
	if len(id) == 0 {
		return responses.BadRequest(fiberCtx, "Invalid params {id}")
	}

	err = company_questionnaire_adp.DeleteQuestionnaire(
		fiberCtx.UserContext(),
		int64(idInt),
		user.CompanyID,
	)
	if errors.Is(err, company_questionnaire_adp.ErrQuestionnaireNotFound) {
		return responses.NotFound(fiberCtx, err.Error())
	}

	if err != nil {
		return responses.InternalServerError(fiberCtx, err)
	}

	return responses.Success(fiberCtx, nil)
}
