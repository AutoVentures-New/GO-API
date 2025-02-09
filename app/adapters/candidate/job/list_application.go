package job

import (
	"context"

	"github.com/hubjob/api/database"
	"github.com/hubjob/api/model"
	"github.com/sirupsen/logrus"
)

func ListJobApplications(
	ctx context.Context,
	candidateID int64,
) ([]model.Application, error) {
	applications := make([]model.Application, 0)

	rows, err := database.Database.QueryContext(
		ctx,
		`SELECT id,company_id,job_id,candidate_id,current_step,status,created_at,updated_at 
				FROM job_applications WHERE candidate_id = ?
				ORDER BY id DESC`,
		candidateID,
	)
	if err != nil {
		logrus.WithError(err).Error("Error to list job application")

		return applications, err
	}

	defer rows.Close()

	for rows.Next() {
		var application model.Application

		err = rows.Scan(
			&application.ID,
			&application.CompanyID,
			&application.JobID,
			&application.CandidateID,
			&application.CurrentStep,
			&application.Status,
			&application.CreatedAt,
			&application.UpdatedAt,
		)
		if err != nil {
			logrus.WithError(err).Error("Error to unmarshal application")

			return applications, err
		}

		applications = append(applications, application)
	}

	return applications, nil
}
