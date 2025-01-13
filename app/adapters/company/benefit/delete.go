package benefit

import (
	"context"

	"github.com/hubjob/api/database"
	"github.com/sirupsen/logrus"
)

func DeleteBenefit(
	ctx context.Context,
	id int64,
	companyID int64,
) error {
	_, err := database.Database.ExecContext(
		ctx,
		`DELETE FROM benefits WHERE id = ? AND company_id = ?`,
		id,
		companyID,
	)
	if err != nil {
		logrus.WithError(err).Error("Error to delete benefit")

		return err
	}

	return nil
}
