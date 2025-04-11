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

		application.Job, err = getJobForApplication(ctx, application.JobID)
		if err != nil {
			return applications, err
		}

		applications = append(applications, application)
	}

	return applications, nil
}

func getJobForApplication(
	ctx context.Context,
	id int64,
) (*model.Job, error) {
	job := model.Job{}

	err := database.Database.QueryRowContext(
		ctx,
		`SELECT id,title,company_id,area_id,is_talent_bank,is_special_needs,description,job_mode,contracting_modality,state,city,responsibilities,questionnaire,video_link,status,publish_at,finish_at,created_at,updated_at 
				FROM jobs WHERE id = ?`,
		id,
	).Scan(
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
		logrus.WithError(err).Error("Error to get job")

		return nil, err
	}

	return &job, nil
}
