package questionnaire_question

import (
	"context"
	"database/sql"
	"errors"

	"github.com/hubjob/api/database"
	"github.com/sirupsen/logrus"
)

var ErrQuestionNotFound = errors.New("Question not found")

func DeleteQuestion(
	ctx context.Context,
	id int64,
	questionnaireID int64,
) error {
	dbTransaction, err := database.Database.Begin()
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		logrus.WithError(err).Error("Error to open transaction")

		return err
	}

	_, err = dbTransaction.ExecContext(
		ctx,
		`DELETE FROM questionnaire_question_answers WHERE questionnaire_id = ? AND questionnaire_question_id = ?`,
		questionnaireID,
		id,
	)
	if err != nil {
		_ = dbTransaction.Rollback()

		logrus.WithError(err).Error("Error to delete answers")

		return err
	}

	result, err := dbTransaction.ExecContext(
		ctx,
		`DELETE FROM questionnaire_questions WHERE questionnaire_id = ? AND id = ?`,
		questionnaireID,
		id,
	)
	if err != nil {
		_ = dbTransaction.Rollback()

		logrus.WithError(err).Error("Error to delete question")

		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		_ = dbTransaction.Rollback()

		logrus.WithError(err).Error("Error to get rows affected")

		return err
	}

	if rowsAffected == 0 {
		_ = dbTransaction.Rollback()

		return ErrQuestionNotFound
	}

	if err := dbTransaction.Commit(); err != nil {
		_ = dbTransaction.Rollback()

		logrus.WithError(err).Error("Error to commit transaction")

		return err
	}

	return nil
}
