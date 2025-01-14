package job

import (
	"context"
	"database/sql"
	"errors"

	"github.com/hubjob/api/database"
	"github.com/hubjob/api/model"
	"github.com/sirupsen/logrus"
)

var ErrJobNotFound = errors.New("Job not found")

func GetJob(
	ctx context.Context,
	id int64,
	companyID int64,
) (model.Job, error) {
	job := model.Job{}

	err := database.Database.QueryRowContext(
		ctx,
		`SELECT id,title,company_id,is_talent_bank,is_special_needs,description,job_mode,contracting_modality,state,city,responsibilities,questionnaire,video_link,status,publish_at,finish_at,created_at,updated_at 
				FROM jobs WHERE company_id = ? AND id = ?`,
		companyID,
		id,
	).Scan(
		&job.ID,
		&job.Title,
		&job.CompanyID,
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
	if errors.Is(err, sql.ErrNoRows) {
		return job, ErrJobNotFound
	}

	if err != nil {
		logrus.WithError(err).Error("Error to get job")

		return job, err
	}

	return job, nil
}
