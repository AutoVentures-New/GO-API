package questionnaire

import (
	"context"
	"fmt"

	questionnaire_question "github.com/hubjob/api/app/adapters/company/questionnaire/question"
	"github.com/hubjob/api/database"
	"github.com/hubjob/api/model"
	"github.com/sirupsen/logrus"
)

func ListQuestionnaires(
	ctx context.Context,
	companyID int64,
) ([]model.Questionnaire, error) {
	questionnaires := make([]model.Questionnaire, 0)

	rows, err := database.Database.QueryContext(
		ctx,
		`SELECT id, name, company_id, created_at, updated_at FROM questionnaires WHERE company_id = ?
				ORDER BY created_at DESC`,
		companyID,
	)
	if err != nil {
		logrus.WithError(err).Error("Error to list questionnaires")

		return nil, err
	}

	defer rows.Close()

	questionnaireIDs := make([]string, 0)

	for rows.Next() {
		questionnaire := model.Questionnaire{}
		err := rows.Scan(
			&questionnaire.ID,
			&questionnaire.Name,
			&questionnaire.CompanyID,
			&questionnaire.CreatedAt,
			&questionnaire.UpdatedAt,
		)
		if err != nil {
			logrus.WithError(err).Error("Error to scan questionnaires")

			return nil, err
		}

		questionnaireIDs = append(questionnaireIDs, fmt.Sprintf("%d", questionnaire.ID))
		questionnaires = append(questionnaires, questionnaire)
	}

	if len(questionnaires) == 0 {
		return questionnaires, nil
	}

	questionsItems, err := questionnaire_question.ListQuestions(ctx, 0, questionnaireIDs)
	if err != nil {
		logrus.WithError(err).Error("Error to list questions")

		return nil, err
	}

	questions := make(map[int64][]model.Question)

	for _, question := range questionsItems {
		if _, ok := questions[question.QuestionnaireID]; !ok {
			questions[question.QuestionnaireID] = make([]model.Question, 0)
		}

		questions[question.QuestionnaireID] = append(questions[question.QuestionnaireID], question)
	}

	for index, questionnaire := range questionnaires {
		if _, ok := questions[questionnaire.ID]; !ok {
			continue
		}

		questionnaires[index].Questions = questions[questionnaire.ID]
	}

	return questionnaires, nil
}
