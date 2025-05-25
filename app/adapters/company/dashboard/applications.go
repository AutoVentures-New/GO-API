package dashboard_adp

import (
	"context"
	"time"

	"github.com/hubjob/api/database"
	"github.com/hubjob/api/model"
	"github.com/sirupsen/logrus"
)

func Applications(
	ctx context.Context,
	companyID int64,
) (model.Applications, error) {
	applications := model.Applications{
		Applications:         make([]model.ApplicationDash, 0),
		ApplicationDashDates: make(map[string]int64),
	}

	rows, err := database.Database.QueryContext(
		ctx,
		`SELECT ap.id,ap.status,ap.job_id,ap.candidate_id,c.name,ap.created_at
				FROM job_applications ap
				JOIN candidates c ON ap.candidate_id = c.id
				WHERE ap.company_id = ? AND ap.created_at > ?
				ORDER BY ap.id DESC`,
		companyID,
		time.Now().UTC().AddDate(0, 0, -90),
	)
	if err != nil {
		logrus.WithError(err).Error("Error to list applications")

		return applications, err
	}

	defer rows.Close()

	for rows.Next() {
		application := model.ApplicationDash{}

		err := rows.Scan(
			&application.ID,
			&application.Status,
			&application.JobID,
			&application.CandidateID,
			&application.CandidateName,
			&application.CreatedAt,
		)
		if err != nil {
			logrus.WithError(err).Error("Error to scan application")

			return applications, err
		}

		applications.Count++

		if application.Status == model.WAITING_EVALUATION {
			applications.WaitingEvaluation++
		}

		if len(applications.Applications) <= 10 {
			applications.Applications = append(applications.Applications, application)
		}

		date := application.CreatedAt.Format("2006-01-02")

		if _, ok := applications.ApplicationDashDates[date]; !ok {
			applications.ApplicationDashDates[date] = 0
		}

		applications.ApplicationDashDates[date]++
	}

	return applications, nil
}
