package questionnaire_question

import (
	"context"

	"github.com/hubjob/api/database"
	"github.com/hubjob/api/model"
	"github.com/sirupsen/logrus"
)

func ListQuestions(
	ctx context.Context,
	questionnaireID int64,
) ([]model.Question, error) {
	questions := make([]model.Question, 0)

	rows, err := database.Database.QueryContext(
		ctx,
		`SELECT id, title, type, questionnaire_id, created_at, updated_at FROM questionnaire_questions WHERE questionnaire_id = ?`,
		questionnaireID,
	)
	if err != nil {
		logrus.WithError(err).Error("Error to list questions")

		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		question := model.Question{}
		err := rows.Scan(
			&question.ID,
			&question.Title,
			&question.Type,
			&question.QuestionnaireID,
			&question.CreatedAt,
			&question.UpdatedAt,
		)
		if err != nil {
			logrus.WithError(err).Error("Error to scan question")

			return nil, err
		}

		answers, err := listAnswers(ctx, question.ID)
		if err != nil {
			return nil, err
		}

		question.Answers = answers

		questions = append(questions, question)
	}

	return questions, nil
}

func listAnswers(
	ctx context.Context,
	questionID int64,
) ([]model.Answer, error) {
	answers := make([]model.Answer, 0)

	rows, err := database.Database.QueryContext(
		ctx,
		`SELECT id, title, questionnaire_id, questionnaire_question_id, is_correct, created_at, updated_at FROM questionnaire_question_answers WHERE questionnaire_question_id = ?`,
		questionID,
	)
	if err != nil {
		logrus.WithError(err).Error("Error to list answers")

		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		answer := model.Answer{}
		err := rows.Scan(
			&answer.ID,
			&answer.Title,
			&answer.QuestionnaireID,
			&answer.QuestionID,
			&answer.IsCorrect,
			&answer.CreatedAt,
			&answer.UpdatedAt,
		)
		if err != nil {
			logrus.WithError(err).Error("Error to scan answer")

			return nil, err
		}

		answers = append(answers, answer)
	}

	return answers, nil
}
