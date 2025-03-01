package job

import (
	"context"

	"github.com/hubjob/api/database"
	"github.com/hubjob/api/model"
	"github.com/sirupsen/logrus"
)

func ListJobApplications(
	ctx context.Context,
	companyID int64,
	jobID int64,
) ([]model.Application, error) {
	applications := make([]model.Application, 0)

	rows, err := database.Database.QueryContext(
		ctx,
		`SELECT j.id,j.company_id,j.job_id,j.candidate_id,j.current_step,j.status,j.created_at,j.updated_at,c.name 
				FROM job_applications j 
				JOIN candidates c on c.id = j.candidate_id
				WHERE j.company_id = ? AND j.job_id = ?
				ORDER BY id DESC`,
		companyID,
		jobID,
	)
	if err != nil {
		logrus.WithError(err).Error("Error to list job application")

		return applications, err
	}

	defer rows.Close()

	for rows.Next() {
		var application model.Application

		var candidate model.Candidate

		err = rows.Scan(
			&application.ID,
			&application.CompanyID,
			&application.JobID,
			&application.CandidateID,
			&application.CurrentStep,
			&application.Status,
			&application.CreatedAt,
			&application.UpdatedAt,

			&candidate.Name,
		)
		if err != nil {
			logrus.WithError(err).Error("Error to unmarshal application")

			return applications, err
		}

		application.Candidate = &candidate
		applications = append(applications, application)
	}

	return applications, nil
}
