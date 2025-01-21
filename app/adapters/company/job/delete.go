package job

import (
	"context"

	"github.com/hubjob/api/database"
	"github.com/sirupsen/logrus"
)

func DeleteJob(
	ctx context.Context,
	id int64,
	companyID int64,
) error {
	dbTransaction, err := database.Database.Begin()
	if err != nil {
		logrus.WithError(err).Error("Error to open db transaction")

		return err
	}

	_, err = dbTransaction.ExecContext(
		ctx,
		`DELETE FROM job_benefits WHERE company_id = ? AND job_id = ?`,
		companyID,
		id,
	)
	if err != nil {
		_ = dbTransaction.Rollback()

		logrus.WithError(err).Error("Error to delete job benefits")

		return err
	}

	_, err = dbTransaction.ExecContext(
		ctx,
		`DELETE FROM job_cultural_fit WHERE company_id = ? AND job_id = ?`,
		companyID,
		id,
	)
	if err != nil {
		_ = dbTransaction.Rollback()

		logrus.WithError(err).Error("Error to delete job cultural fit")

		return err
	}

	_, err = dbTransaction.ExecContext(
		ctx,
		`DELETE FROM job_requirements WHERE company_id = ? AND job_id = ?`,
		companyID,
		id,
	)
	if err != nil {
		_ = dbTransaction.Rollback()

		logrus.WithError(err).Error("Error to delete job requirements")

		return err
	}

	result, err := dbTransaction.ExecContext(
		ctx,
		`DELETE FROM jobs WHERE id = ? AND company_id = ?`,
		id,
		companyID,
	)
	if err != nil {
		_ = dbTransaction.Rollback()

		logrus.WithError(err).Error("Error to delete job")

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

		return ErrJobNotFound
	}

	if err := dbTransaction.Commit(); err != nil {
		_ = dbTransaction.Rollback()

		logrus.WithError(err).Error("Error to commit transaction")

		return err
	}

	return nil
}
