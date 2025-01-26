package get_application

import (
	"context"

	"github.com/hubjob/api/database"
	"github.com/hubjob/api/model"
	"github.com/sirupsen/logrus"
)

func ListQuestions(
	ctx context.Context,
	jobID int64,
) ([]model.ApplicationQuestion, error) {
	questions := make([]model.ApplicationQuestion, 0)

	rows, err := database.Database.QueryContext(
		ctx,
		`SELECT 
    				questionnaire_questions.id, 
    				questionnaire_questions.title, 
    				questionnaire_questions.type 
				FROM job_questions
				JOIN questionnaire_questions ON job_questions.question_id = questionnaire_questions.id 
				WHERE job_questions.job_id = ?`,
		jobID,
	)
	if err != nil {
		logrus.WithError(err).Error("Error to list questions")

		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		question := model.ApplicationQuestion{}
		err := rows.Scan(
			&question.ID,
			&question.Title,
			&question.Type,
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
) ([]model.ApplicationAnswer, error) {
	answers := make([]model.ApplicationAnswer, 0)

	rows, err := database.Database.QueryContext(
		ctx,
		`SELECT id, title FROM questionnaire_question_answers WHERE questionnaire_question_id = ?`,
		questionID,
	)
	if err != nil {
		logrus.WithError(err).Error("Error to list answers")

		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		answer := model.ApplicationAnswer{}
		err := rows.Scan(
			&answer.ID,
			&answer.Title,
		)
		if err != nil {
			logrus.WithError(err).Error("Error to scan answer")

			return nil, err
		}

		answers = append(answers, answer)
	}

	return answers, nil
}
