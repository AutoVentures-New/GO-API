package configuration

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/hubjob/api/app/adapters/area"
	company_job_adp "github.com/hubjob/api/app/adapters/company/job"
	"github.com/hubjob/api/handler/responses"
	"github.com/hubjob/api/model"
	"os"
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

func ListStateCities(fiberCtx *fiber.Ctx) error {
	idInt, _ := strconv.Atoi(fiberCtx.Query("company_id", "0"))

	jobs, err := company_job_adp.ListStateCities(
		fiberCtx.UserContext(),
		int64(idInt),
	)
	if err != nil {
		return responses.InternalServerError(fiberCtx, err)
	}

	return responses.Success(fiberCtx, jobs)
}

func ListExampleQuestions(fiberCtx *fiber.Ctx) error {
	questions, err := loadQuestionsFromFile()
	if err != nil {
		return responses.InternalServerError(fiberCtx, err)
	}

	return responses.Success(fiberCtx, questions)
}

func loadQuestionsFromFile() ([]model.ExampleQuestions, error) {
	file, err := os.ReadFile("./storage/example_questions.json")
	if err != nil {
		return nil, err
	}

	var questions []model.ExampleQuestions
	err = json.Unmarshal(file, &questions)
	if err != nil {
		return nil, err
	}

	return questions, nil
}
