package job

import (
	"context"
	"time"

	"github.com/hubjob/api/database"
	"github.com/hubjob/api/model"
	"github.com/sirupsen/logrus"
)

func CanceledApplication(
	ctx context.Context,
	application model.Application,
) (model.Application, error) {
	dbTransaction, err := database.Database.Begin()
	if err != nil {
		logrus.WithError(err).Error("Error to open db transaction")

		return application, err
	}

	application.Status = model.CANCELED
	application.UpdatedAt = time.Now().UTC()

	_, err = dbTransaction.ExecContext(
		ctx,
		`UPDATE job_applications set status = ?, updated_at = ? WHERE id = ?`,
		application.Status,
		application.UpdatedAt,
		application.ID,
	)
	if err != nil {
		logrus.WithError(err).Error("Error to update job application")

		return application, err
	}

	if err := dbTransaction.Commit(); err != nil {
		_ = dbTransaction.Rollback()

		logrus.WithError(err).Error("Error to commit transaction")

		return application, err
	}

	return application, nil
}
