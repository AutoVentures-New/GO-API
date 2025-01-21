package job

import (
	"context"
	"database/sql"
	"encoding/json"
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

	dbTransaction, err := database.Database.Begin()
	if err != nil {
		logrus.WithError(err).Error("Error to open db transaction")

		return job, err
	}

	result, err := dbTransaction.ExecContext(
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

		_ = dbTransaction.Rollback()

		return job, err
	}

	job.ID, err = result.LastInsertId()
	if err != nil {
		logrus.WithError(err).Error("Error to get last insert job id")

		_ = dbTransaction.Rollback()

		return job, err
	}

	job.JobCulturalFit.CompanyID = job.CompanyID
	job.JobCulturalFit.JobID = job.ID
	job.JobCulturalFit.CreatedAt = job.CreatedAt
	job.JobCulturalFit.UpdatedAt = job.CreatedAt

	answersString, err := json.Marshal(job.JobCulturalFit.Answers)
	if err != nil {
		logrus.WithError(err).Error("Error to marshal job cultural fit answers")

		_ = dbTransaction.Rollback()

		return job, err
	}

	resultJobCulturalFit, err := dbTransaction.ExecContext(
		ctx,
		`INSERT INTO job_cultural_fit(company_id,job_id,answers,created_at,updated_at) 
					VALUES(?,?,?,?,?)`,
		job.JobCulturalFit.CompanyID,
		job.JobCulturalFit.JobID,
		answersString,
		job.JobCulturalFit.CreatedAt,
		job.JobCulturalFit.UpdatedAt,
	)
	if err != nil {
		logrus.WithError(err).Error("Error to insert job cultural fit")

		_ = dbTransaction.Rollback()

		return job, err
	}

	job.JobCulturalFit.ID, err = resultJobCulturalFit.LastInsertId()
	if err != nil {
		logrus.WithError(err).Error("Error to get last insert job cultural fit id")

		_ = dbTransaction.Rollback()

		return job, err
	}

	if err := CreateJobRequirement(ctx, dbTransaction, &job); err != nil {
		_ = dbTransaction.Rollback()

		return job, err
	}

	if err := CreateJobBenefits(ctx, dbTransaction, &job); err != nil {
		_ = dbTransaction.Rollback()

		return job, err
	}

	if err := dbTransaction.Commit(); err != nil {
		_ = dbTransaction.Rollback()

		logrus.WithError(err).Error("Error to commit transaction")

		return job, err
	}

	return job, nil
}

func CreateJobRequirement(
	ctx context.Context,
	dbTransaction *sql.Tx,
	job *model.Job,
) error {
	job.JobRequirement.CompanyID = job.CompanyID
	job.JobRequirement.JobID = job.ID
	job.JobRequirement.CreatedAt = job.CreatedAt
	job.JobRequirement.UpdatedAt = job.CreatedAt

	itemsString, err := json.Marshal(job.JobRequirement.Items)
	if err != nil {
		logrus.WithError(err).Error("Error to marshal job requirements items")

		return err
	}

	resultJobRequirement, err := dbTransaction.ExecContext(
		ctx,
		`INSERT INTO job_requirements(company_id,job_id,items,min_match,created_at,updated_at) 
					VALUES(?,?,?,?,?,?)`,
		job.JobRequirement.CompanyID,
		job.JobRequirement.JobID,
		itemsString,
		job.JobRequirement.MinMatch,
		job.JobRequirement.CreatedAt,
		job.JobRequirement.UpdatedAt,
	)
	if err != nil {
		logrus.WithError(err).Error("Error to insert job requirements")

		return err
	}

	job.JobRequirement.ID, err = resultJobRequirement.LastInsertId()
	if err != nil {
		logrus.WithError(err).Error("Error to get last insert job requirements id")

		return err
	}

	return nil
}

func CreateJobBenefits(
	ctx context.Context,
	dbTransaction *sql.Tx,
	job *model.Job,
) error {
	for _, benefit := range job.Benefits {
		_, err := dbTransaction.ExecContext(
			ctx,
			`INSERT INTO job_benefits(company_id,job_id,benefit_id,created_at,updated_at) 
					VALUES(?,?,?,?,?)`,
			job.CompanyID,
			job.ID,
			benefit.ID,
			job.CreatedAt,
			job.UpdatedAt,
		)
		if err != nil {
			logrus.WithError(err).Error("Error to insert job benefit")

			return err
		}
	}

	return nil
}
