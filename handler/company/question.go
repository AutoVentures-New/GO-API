package company

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	company_questionnaire_question_adp "github.com/hubjob/api/app/adapters/company/questionnaire/question"
	"github.com/hubjob/api/handler/responses"
	"github.com/hubjob/api/model"
)

type CreateQuestionRequest struct {
	Title   string `json:"title"`
	Type    string `json:"type"`
	Answers []struct {
		Title     string `json:"title"`
		IsCorrect bool   `json:"is_correct"`
	} `json:"answers"`
}

func CreateQuestion(fiberCtx *fiber.Ctx) error {
	//user := fiberCtx.Locals("user").(model.User)
	request := new(CreateQuestionRequest)

	if err := fiberCtx.BodyParser(&request); err != nil {
		return responses.InvalidBodyRequest(fiberCtx, err)
	}

	questionnaireId := fiberCtx.Params("questionnaire_id")
	if len(questionnaireId) == 0 {
		return responses.BadRequest(fiberCtx, "Params {questionnaire_id} is required")
	}

	questionnaireIdInt, err := strconv.Atoi(questionnaireId)
	if err != nil {
		return responses.BadRequest(fiberCtx, "Invalid params {questionnaire_id}")
	}

	answers := make([]model.Answer, 0)

	for _, answer := range request.Answers {
		answers = append(answers, model.Answer{
			Title:     answer.Title,
			IsCorrect: answer.IsCorrect,
		})
	}

	question, err := company_questionnaire_question_adp.CreateQuestion(
		fiberCtx.UserContext(),
		int64(questionnaireIdInt),
		model.Question{
			Title:   request.Title,
			Type:    request.Type,
			Answers: answers,
		},
	)
	if errors.Is(err, company_questionnaire_question_adp.ErrQuestionAlreadyExists) {
		return responses.Conflict(fiberCtx, "COMPANY|QUESTIONNAIRE_QUESTION|ALREADY_EXISTS", "Question already exists")
	}

	if err != nil {
		return responses.InternalServerError(fiberCtx, err)
	}

	return responses.Success(fiberCtx, question)
}

func ListQuestions(fiberCtx *fiber.Ctx) error {
	//user := fiberCtx.Locals("user").(model.User)

	questionnaireId := fiberCtx.Params("questionnaire_id")
	if len(questionnaireId) == 0 {
		return responses.BadRequest(fiberCtx, "Params {questionnaire_id} is required")
	}

	questionnaireIdInt, err := strconv.Atoi(questionnaireId)
	if err != nil {
		return responses.BadRequest(fiberCtx, "Invalid params {questionnaire_id}")
	}

	questions, err := company_questionnaire_question_adp.ListQuestions(
		fiberCtx.UserContext(),
		int64(questionnaireIdInt),
	)
	if err != nil {
		return responses.InternalServerError(fiberCtx, err)
	}

	return responses.Success(fiberCtx, questions)
}

type UpdateQuestionRequest struct {
	Title   string `json:"title"`
	Type    string `json:"type"`
	Answers []struct {
		ID        int64  `json:"id"`
		Title     string `json:"title"`
		IsCorrect bool   `json:"is_correct"`
	} `json:"answers"`
}

func UpdateQuestion(fiberCtx *fiber.Ctx) error {
	//user := fiberCtx.Locals("user").(model.User)
	request := new(UpdateQuestionRequest)

	if err := fiberCtx.BodyParser(&request); err != nil {
		return responses.InvalidBodyRequest(fiberCtx, err)
	}

	questionnaireId := fiberCtx.Params("questionnaire_id")
	if len(questionnaireId) == 0 {
		return responses.BadRequest(fiberCtx, "Params {questionnaire_id} is required")
	}

	questionnaireIdInt, err := strconv.Atoi(questionnaireId)
	if err != nil {
		return responses.BadRequest(fiberCtx, "Invalid params {questionnaire_id}")
	}

	id := fiberCtx.Params("id")
	if len(id) == 0 {
		return responses.BadRequest(fiberCtx, "Params {id} is required")
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		return responses.BadRequest(fiberCtx, "Invalid params {id}")
	}

	answers := make([]model.Answer, 0)

	for _, answer := range request.Answers {
		answers = append(answers, model.Answer{
			ID:        answer.ID,
			Title:     answer.Title,
			IsCorrect: answer.IsCorrect,
		})
	}

	question, err := company_questionnaire_question_adp.UpdateQuestion(
		fiberCtx.UserContext(),
		model.Question{
			ID:              int64(idInt),
			Title:           request.Title,
			Type:            request.Type,
			QuestionnaireID: int64(questionnaireIdInt),
			Answers:         answers,
		},
	)
	if err != nil {
		return responses.InternalServerError(fiberCtx, err)
	}

	return responses.Success(fiberCtx, question)
}

func DeleteQuestion(fiberCtx *fiber.Ctx) error {
	//user := fiberCtx.Locals("user").(model.User)

	questionnaireId := fiberCtx.Params("questionnaire_id")
	if len(questionnaireId) == 0 {
		return responses.BadRequest(fiberCtx, "Params {questionnaire_id} is required")
	}

	questionnaireIdInt, err := strconv.Atoi(questionnaireId)
	if err != nil {
		return responses.BadRequest(fiberCtx, "Invalid params {questionnaire_id}")
	}

	id := fiberCtx.Params("id")
	if len(id) == 0 {
		return responses.BadRequest(fiberCtx, "Params {id} is required")
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		return responses.BadRequest(fiberCtx, "Invalid params {id}")
	}

	err = company_questionnaire_question_adp.DeleteQuestion(
		fiberCtx.UserContext(),
		int64(idInt),
		int64(questionnaireIdInt),
	)
	if errors.Is(err, company_questionnaire_question_adp.ErrQuestionNotFound) {
		return responses.NotFound(fiberCtx, err.Error())
	}

	if err != nil {
		return responses.InternalServerError(fiberCtx, err)
	}

	return responses.Success(fiberCtx, nil)
}
