package questionnaire_adp

import (
	"context"
	"time"

	"github.com/hubjob/api/database"
	"github.com/hubjob/api/model"
	"github.com/sirupsen/logrus"
)

func UpdateCandidateQuestionnaire(
	ctx context.Context,
	questionnaire model.CandidateQuestionnaire,
) error {
	dbTransaction, err := database.Database.Begin()
	if err != nil {
		logrus.WithError(err).Error("Error to open db transaction")

		return err
	}

	questionnaire.UpdatedAt = time.Now().UTC()

	_, err = dbTransaction.ExecContext(
		ctx,
		`UPDATE candidate_questionnaires set
    				bucket_name = ?, 
    				result_file_path = ?,
    				updated_at = ?
				WHERE id = ?`,
		questionnaire.BucketName,
		questionnaire.ResultFilePath,
		questionnaire.UpdatedAt,
		questionnaire.ID,
	)
	if err != nil {
		_ = dbTransaction.Rollback()

		logrus.WithError(err).Error("Error to update candidate questionnaire")

		return err
	}

	err = dbTransaction.Commit()
	if err != nil {
		_ = dbTransaction.Rollback()

		logrus.WithError(err).Error("Error to commit db transaction")

		return err
	}

	return nil
}
