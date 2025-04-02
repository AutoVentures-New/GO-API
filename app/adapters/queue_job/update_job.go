package queue_job

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hubjob/api/database"
	"github.com/hubjob/api/model"
	"github.com/sirupsen/logrus"
)

func UpdateJobsToProcessing(
	ctx context.Context,
	jobID []string,
) error {
	_, err := database.Database.ExecContext(
		ctx,
		fmt.Sprintf(`UPDATE queue_jobs set status='PROCESSING',updated_at = ? WHERE id IN (%s)`, strings.Join(jobID, ",")),
		time.Now().UTC(),
	)
	if err != nil {
		logrus.WithError(err).Error("Error to update jobs to processing")

		return err
	}

	return nil
}

func UpdateJob(
	ctx context.Context,
	job model.QueueJob,
) error {
	_, err := database.Database.ExecContext(
		ctx,
		`UPDATE queue_jobs set status=?,updated_at = ? WHERE id = ?`,
		job.Status,
		time.Now().UTC(),
		job.ID,
	)
	if err != nil {
		logrus.WithError(err).Error("Error to update job")

		return err
	}

	return nil
}
