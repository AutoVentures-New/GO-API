package questionnaire

import (
	"context"

	"github.com/hubjob/api/database"
	"github.com/sirupsen/logrus"
)

func DeleteQuestionnaire(
	ctx context.Context,
	id int64,
	companyID int64,
) error {
	result, err := database.Database.ExecContext(
		ctx,
		`DELETE FROM questionnaires WHERE id = ? AND company_id = ?`,
		id,
		companyID,
	)
	if err != nil {
		logrus.WithError(err).Error("Error to delete questionnaire")

		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logrus.WithError(err).Error("Error to get rows affected")

		return err
	}

	if rowsAffected == 0 {
		return ErrQuestionnaireNotFound
	}

	return nil
}
