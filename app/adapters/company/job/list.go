package job

import (
	"context"

	"github.com/hubjob/api/database"
	"github.com/hubjob/api/model"
	"github.com/sirupsen/logrus"
)

func ListJobs(
	ctx context.Context,
	companyID int64,
) ([]model.Job, error) {
	jobs := make([]model.Job, 0)

	rows, err := database.Database.QueryContext(
		ctx,
		`SELECT id,title,company_id,area_id,is_talent_bank,is_special_needs,description,job_mode,contracting_modality,state,city,responsibilities,questionnaire,video_link,status,publish_at,finish_at,created_at,updated_at
				FROM jobs WHERE company_id = ?`,
		companyID,
	)
	if err != nil {
		logrus.WithError(err).Error("Error to list jobs")

		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		job := model.Job{}
		err := rows.Scan(
			&job.ID,
			&job.Title,
			&job.CompanyID,
			&job.AreaID,
			&job.IsTalentBank,
			&job.IsSpecialNeeds,
			&job.Description,
			&job.JobMode,
			&job.ContractingModality,
			&job.State,
			&job.City,
			&job.Responsibilities,
			&job.Questionnaire,
			&job.VideoLink,
			&job.Status,
			&job.PublishAt,
			&job.FinishAt,
			&job.CreatedAt,
			&job.UpdatedAt,
		)
		if err != nil {
			logrus.WithError(err).Error("Error to scan jobs")

			return nil, err
		}

		jobs = append(jobs, job)
	}

	return jobs, nil
}
