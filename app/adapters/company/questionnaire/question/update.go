package questionnaire_question

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hubjob/api/database"
	"github.com/hubjob/api/model"
	"github.com/sirupsen/logrus"
)

func UpdateQuestion(
	ctx context.Context,
	question model.Question,
) (model.Question, error) {
	dbTransaction, err := database.Database.Begin()
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		logrus.WithError(err).Error("Error to open transaction")

		return question, err
	}

	question.UpdatedAt = time.Now().UTC()

	_, err = dbTransaction.ExecContext(
		ctx,
		`UPDATE questionnaire_questions set title = ?, type = ?, updated_at = ? WHERE id = ?`,
		question.Title,
		question.Type,
		question.UpdatedAt,
		question.ID,
	)
	if err != nil {
		_ = dbTransaction.Rollback()

		logrus.WithError(err).Error("Error to update question")

		return question, err
	}

	if question.Type != model.OPEN_FIELD {
		newAnswer := make([]model.Answer, 0)
		idsToIgnore := make([]string, 0)

		for _, answer := range question.Answers {
			if answer.ID == 0 {
				answer.QuestionnaireID = question.QuestionnaireID
				answer.QuestionID = question.ID
				answer.CreatedAt = question.UpdatedAt
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
			} else {
				answer.UpdatedAt = question.UpdatedAt

				err = updateAnswer(ctx, dbTransaction, answer)
				if err != nil {
					_ = dbTransaction.Rollback()

					return question, err
				}
			}

			idsToIgnore = append(idsToIgnore, fmt.Sprintf("%d", answer.ID))
			newAnswer = append(newAnswer, answer)
		}

		err = deleteAnswers(ctx, dbTransaction, question.ID, idsToIgnore)
		if err != nil {
			_ = dbTransaction.Rollback()

			return question, err
		}

		question.Answers = newAnswer
	}

	if err := dbTransaction.Commit(); err != nil {
		_ = dbTransaction.Rollback()

		logrus.WithError(err).Error("Error to commit transaction")

		return question, err
	}

	return question, nil
}

func updateAnswer(
	ctx context.Context,
	dbTransaction *sql.Tx,
	answer model.Answer,
) error {
	_, err := dbTransaction.ExecContext(
		ctx,
		`UPDATE questionnaire_question_answers set title = ?, is_correct = ?, updated_at = ? WHERE id = ?`,
		answer.Title,
		answer.IsCorrect,
		answer.UpdatedAt,
		answer.ID,
	)
	if err != nil {
		logrus.WithError(err).Error("Error to update answer")

		return err
	}

	return nil
}

func deleteAnswers(
	ctx context.Context,
	dbTransaction *sql.Tx,
	questionID int64,
	idsToIgnore []string,
) error {
	query := fmt.Sprintf(
		`DELETE FROM questionnaire_question_answers WHERE questionnaire_question_id = ? AND id NOT IN(%s)`,
		strings.Join(idsToIgnore, ","),
	)

	_, err := dbTransaction.ExecContext(
		ctx,
		query,
		questionID,
	)
	if err != nil {
		logrus.WithError(err).Error("Error to delete answers")

		return err
	}

	return nil
}
