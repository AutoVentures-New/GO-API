package job

import (
	"context"
	"time"

	"github.com/hubjob/api/database"
	"github.com/sirupsen/logrus"
)

func FinishJob(
	ctx context.Context,
	jobID int64,
	companyID int64,
) error {
	dbTransaction, err := database.Database.Begin()
	if err != nil {
		logrus.WithError(err).Error("Error to open db transaction")

		return err
	}

	_, err = dbTransaction.ExecContext(
		ctx,
		`UPDATE jobs set finish_at = ?,status = ?,updated_at = ? WHERE id = ? AND company_id = ? and status != 'FINISHED'`,
		time.Now().UTC(),
		"FINISHED",
		time.Now().UTC(),
		jobID,
		companyID,
	)
	if err != nil {
		logrus.WithError(err).Error("Error to update job to finish")

		_ = dbTransaction.Rollback()

		return err
	}

	if err := dbTransaction.Commit(); err != nil {
		_ = dbTransaction.Rollback()

		logrus.WithError(err).Error("Error to commit transaction")

		return err
	}

	return nil
}
