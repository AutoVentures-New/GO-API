package steps

import (
	"context"
	"database/sql"
	"time"

	"github.com/hubjob/api/model"
	"github.com/sirupsen/logrus"
)

func updateApplication(
	ctx context.Context,
	dbTransaction *sql.Tx,
	application *model.Application,
	reproved bool,
) error {
	application.UpdatedAt = time.Now().UTC()

	if reproved {
		application.Status = model.REPROVED
	} else {
		for index, value := range application.Steps {
			if value == application.CurrentStep && len(application.Steps) > index+1 {
				application.CurrentStep = application.Steps[index+1]

				break
			}
		}
	}

	_, err := dbTransaction.ExecContext(
		ctx,
		`UPDATE job_applications set current_step = ?, status = ?, updated_at = ? WHERE id = ?`,
		application.CurrentStep,
		application.Status,
		application.UpdatedAt,
		application.ID,
	)
	if err != nil {
		logrus.WithError(err).Error("Error to update job application")

		return err
	}

	return nil
}
