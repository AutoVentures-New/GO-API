package queue_job

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hubjob/api/database"
	"github.com/hubjob/api/model"
	"github.com/sirupsen/logrus"
)

func ListJobs(
	ctx context.Context,
) ([]model.QueueJob, error) {
	jobs := make([]model.QueueJob, 0)
	jobIds := make([]string, 0)

	rows, err := database.Database.QueryContext(
		ctx,
		`SELECT id,type,status,configurations,created_at,updated_at FROM queue_jobs WHERE status = 'PENDING'`,
	)
	if err != nil {
		logrus.WithError(err).Error("Error to list queue jobs")

		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		job := model.QueueJob{}

		var configurations []byte

		err := rows.Scan(
			&job.ID,
			&job.Type,
			&job.Status,
			&configurations,
			&job.CreatedAt,
			&job.UpdatedAt,
		)
		if err != nil {
			logrus.WithError(err).Error("Error to scan queue jobs")

			return nil, err
		}

		err = json.Unmarshal(configurations, &job.Configurations)
		if err != nil {
			logrus.WithError(err).Error("Error to unmarshal jobs configurations")

			return nil, err
		}

		jobIds = append(jobIds, fmt.Sprint(job.ID))
		jobs = append(jobs, job)
	}

	if len(jobIds) == 0 {
		return jobs, nil
	}

	if err = UpdateJobsToProcessing(ctx, jobIds); err != nil {
		return nil, err
	}

	return jobs, nil
}
