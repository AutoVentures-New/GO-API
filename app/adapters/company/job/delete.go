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
	result, err := database.Database.ExecContext(
		ctx,
		`DELETE FROM jobs WHERE id = ? AND company_id = ?`,
		id,
		companyID,
	)
	if err != nil {
		logrus.WithError(err).Error("Error to delete job")

		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logrus.WithError(err).Error("Error to get rows affected")

		return err
	}

	if rowsAffected == 0 {
		return ErrJobNotFound
	}

	return nil
}
