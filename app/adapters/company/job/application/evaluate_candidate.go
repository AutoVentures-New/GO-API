package job

import (
	"context"
	"time"

	"github.com/hubjob/api/database"
	"github.com/sirupsen/logrus"
)

func EvaluateCandidate(
	ctx context.Context,
	applicationID int64,
	companyID int64,
	status string,
) error {
	_, err := database.Database.ExecContext(
		ctx,
		`UPDATE job_applications set status = ?, updated_at = ? WHERE id = ? AND company_id = ?`,
		status,
		time.Now().UTC(),
		applicationID,
		companyID,
	)
	if err != nil {
		logrus.WithError(err).Error("Error to evaluate candidate")

		return err
	}

	return nil
}
