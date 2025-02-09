package profile

import (
	"context"
	"time"

	"github.com/hubjob/api/database"
	"github.com/hubjob/api/model"
	"github.com/sirupsen/logrus"
)

func UpdateEmail(
	ctx context.Context,
	candidate model.Candidate,
	code string,
) (model.Candidate, error) {
	dbTransaction, err := database.Database.Begin()
	if err != nil {
		logrus.WithError(err).Error("Error to open db transaction")

		return model.Candidate{}, err
	}

	candidate.UpdatedAt = time.Now().UTC()

	_, err = dbTransaction.ExecContext(
		ctx,
		`UPDATE candidates set
    				email = ?,
    				updated_at = ?
				WHERE id = ?`,
		candidate.Email,
		candidate.UpdatedAt,
		candidate.ID,
	)
	if err != nil {
		_ = dbTransaction.Rollback()

		logrus.WithError(err).Error("Error to update candidate")

		return model.Candidate{}, err
	}

	_, err = dbTransaction.ExecContext(
		ctx,
		`DELETE FROM email_validations WHERE email = ? AND code = ?`,
		candidate.Email,
		code,
	)
	if err != nil {
		_ = dbTransaction.Rollback()

		logrus.WithError(err).Error("Error to delete email validation")

		return model.Candidate{}, err
	}

	err = dbTransaction.Commit()
	if err != nil {
		_ = dbTransaction.Rollback()

		logrus.WithError(err).Error("Error to commit db transaction")

		return model.Candidate{}, err
	}

	return candidate, nil
}
