package steps

import (
	"context"
	"database/sql"
	"github.com/hubjob/api/model"
	"github.com/sirupsen/logrus"
	"time"
)

func updateApplication(
	ctx context.Context,
	dbTransaction *sql.Tx,
	application *model.Application,
	reproved bool,
) error {
	if reproved {
		application.Status = model.REPROVED
	} else {
		if application.CurrentStep == application.Steps[len(application.Steps)-1] {
			application.Status = model.WAITING_EVALUATION
		} else {
			for index, value := range application.Steps {
				if value == application.CurrentStep && len(application.Steps) > index+1 {
					application.CurrentStep = application.Steps[index+1]

					break
				}
			}
		}
	}

	application.UpdatedAt = time.Now().UTC()

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
