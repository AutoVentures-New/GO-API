package questionnaire_question

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/hubjob/api/database"
	"github.com/hubjob/api/model"
	"github.com/sirupsen/logrus"
)

var ErrQuestionAlreadyExists = errors.New("question already exists")

func CreateQuestion(
	ctx context.Context,
	questionnaireID int64,
	question model.Question,
) (model.Question, error) {
	var count int

	err := database.Database.QueryRowContext(
		ctx,
		`SELECT COUNT(0) FROM questionnaire_questions WHERE questionnaire_id = ? AND title = ?`,
		questionnaireID,
		question.Title,
	).Scan(&count)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		logrus.WithError(err).Error("Error to validate question")

		return question, err
	}

	if count > 0 {
		return question, ErrQuestionAlreadyExists
	}

	dbTransaction, err := database.Database.Begin()
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		logrus.WithError(err).Error("Error to open transaction")

		return question, err
	}

	question.QuestionnaireID = questionnaireID
	question.CreatedAt = time.Now().UTC()
	question.UpdatedAt = question.CreatedAt

	result, err := dbTransaction.ExecContext(
		ctx,
		`INSERT INTO questionnaire_questions(title,type,questionnaire_id,created_at,updated_at) VALUES(?,?,?,?,?)`,
		question.Title,
		question.Type,
		question.QuestionnaireID,
		question.CreatedAt,
		question.UpdatedAt,
	)
	if err != nil {
		_ = dbTransaction.Rollback()

		logrus.WithError(err).Error("Error to insert question")

		return question, err
	}

	question.ID, err = result.LastInsertId()
	if err != nil {
		_ = dbTransaction.Rollback()

		logrus.WithError(err).Error("Error to get last insert questionnaire id")

		return question, err
	}

	newAnswers := make([]model.Answer, 0)

	for _, answer := range question.Answers {
		answer.QuestionnaireID = question.QuestionnaireID
		answer.QuestionID = question.ID
		answer.CreatedAt = question.CreatedAt
		answer.UpdatedAt = answer.CreatedAt

		answer.ID, err = createAnswer(
			ctx,
			dbTransaction,
			answer,
		)
		if err != nil {
			_ = dbTransaction.Rollback()

			return question, err
		}

		newAnswers = append(newAnswers, answer)
	}

	question.Answers = newAnswers

	if err := dbTransaction.Commit(); err != nil {
		_ = dbTransaction.Rollback()

		logrus.WithError(err).Error("Error to commit transaction")

		return question, err
	}

	return question, nil
}

func createAnswer(
	ctx context.Context,
	dbTransaction *sql.Tx,
	answer model.Answer,
) (int64, error) {
	result, err := dbTransaction.ExecContext(
		ctx,
		`INSERT INTO questionnaire_question_answers(title,questionnaire_id,questionnaire_question_id,is_correct,created_at,updated_at) VALUES(?,?,?,?,?,?)`,
		answer.Title,
		answer.QuestionnaireID,
		answer.QuestionID,
		answer.IsCorrect,
		answer.CreatedAt,
		answer.UpdatedAt,
	)
	if err != nil {
		logrus.WithError(err).Error("Error to insert answer")

		return 0, err
	}

	answerId, err := result.LastInsertId()
	if err != nil {
		logrus.WithError(err).Error("Error to get last insert answer id")

		return 0, err
	}

	return answerId, nil
}
