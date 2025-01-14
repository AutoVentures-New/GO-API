package job

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/hubjob/api/database"
	"github.com/hubjob/api/model"
	"github.com/sirupsen/logrus"
)

var ErrJobAlreadyExists = errors.New("job already exists")

func CreateJob(
	ctx context.Context,
	job model.Job,
) (model.Job, error) {
	var count int

	err := database.Database.QueryRowContext(
		ctx,
		`SELECT COUNT(0) FROM jobs WHERE company_id = ? AND title = ? AND status = 'ACTIVE' LIMIT 1`,
		job.CompanyID,
		job.Title,
	).Scan(&count)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		logrus.WithError(err).Error("Error to validate job")

		return job, err
	}

	if count > 0 {
		return job, ErrJobAlreadyExists
	}

	job.Status = "ACTIVE"
	job.CreatedAt = time.Now().UTC()
	job.UpdatedAt = job.CreatedAt

	result, err := database.Database.ExecContext(
		ctx,
		`INSERT INTO jobs(title,company_id,is_talent_bank,is_special_needs,description,job_mode,contracting_modality,state,city,responsibilities,questionnaire,video_link,status,publish_at,finish_at,created_at,updated_at) 
					VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`,
		job.Title,
		job.CompanyID,
		job.IsTalentBank,
		job.IsSpecialNeeds,
		job.Description,
		job.JobMode,
		job.ContractingModality,
		job.State,
		job.City,
		job.Responsibilities,
		job.Questionnaire,
		job.VideoLink,
		job.Status,
		job.PublishAt,
		job.FinishAt,
		job.CreatedAt,
		job.UpdatedAt,
	)
	if err != nil {
		logrus.WithError(err).Error("Error to insert job")

		return job, err
	}

	job.ID, err = result.LastInsertId()
	if err != nil {
		logrus.WithError(err).Error("Error to get last insert job id")

		return job, err
	}

	return job, nil
}
