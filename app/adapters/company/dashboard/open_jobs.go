package dashboard_adp

import (
	"context"
	"time"

	"github.com/hubjob/api/database"
	"github.com/hubjob/api/model"
	"github.com/sirupsen/logrus"
)

func OpenJobs(
	ctx context.Context,
	companyID int64,
) (model.OpenJobs, error) {
	jobs := model.OpenJobs{
		Jobs: make([]struct {
			ID    int64  `json:"id"`
			Title string `json:"title"`
		}, 0),
	}

	rows, err := database.Database.QueryContext(
		ctx,
		`SELECT id,title,finish_at
				FROM jobs WHERE company_id = ? AND status = 'ACTIVE'
				ORDER BY id DESC`,
		companyID,
	)
	if err != nil {
		logrus.WithError(err).Error("Error to list jobs opens")

		return jobs, err
	}

	defer rows.Close()

	fiveDays := time.Now().UTC().AddDate(0, 0, 5)

	for rows.Next() {
		job := model.Job{}

		err := rows.Scan(
			&job.ID,
			&job.Title,
			&job.FinishAt,
		)
		if err != nil {
			logrus.WithError(err).Error("Error to scan jobs open")

			return jobs, err
		}

		jobs.Count++

		if len(jobs.Jobs) <= 10 {
			jobs.Jobs = append(jobs.Jobs, struct {
				ID    int64  `json:"id"`
				Title string `json:"title"`
			}{ID: job.ID, Title: job.Title})
		}

		if job.FinishAt.Before(fiveDays) {
			jobs.CloseToFinish++
		}
	}

	return jobs, nil
}
