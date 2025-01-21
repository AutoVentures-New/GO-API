package job

import (
	"context"
	"database/sql"
	"encoding/json"
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

	job.JobCulturalFit, err = getJobCulturalFit(ctx, id, companyID)
	if err != nil {
		return job, err
	}

	job.JobRequirement, err = getJobRequirements(ctx, id, companyID)
	if err != nil {
		return job, err
	}

	job.Benefits, err = getJobBenefits(ctx, id, companyID)
	if err != nil {
		return job, err
	}

	return job, nil
}

func getJobCulturalFit(
	ctx context.Context,
	id int64,
	companyID int64,
) (model.JobCulturalFit, error) {
	jobCulturalFit := model.JobCulturalFit{}

	var answersJSON []byte

	err := database.Database.QueryRowContext(
		ctx,
		`SELECT id,company_id,job_id,answers,created_at,updated_at 
				FROM job_cultural_fit WHERE company_id = ? AND job_id = ? LIMIT 1`,
		companyID,
		id,
	).Scan(
		&jobCulturalFit.ID,
		&jobCulturalFit.CompanyID,
		&jobCulturalFit.JobID,
		&answersJSON,
		&jobCulturalFit.CreatedAt,
		&jobCulturalFit.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return jobCulturalFit, ErrJobNotFound
	}

	if err != nil {
		logrus.WithError(err).Error("Error to get job cultural fit")

		return jobCulturalFit, err
	}

	err = json.Unmarshal(answersJSON, &jobCulturalFit.Answers)
	if err != nil {
		logrus.WithError(err).Error("Error to unmarshal job cultural fit")

		return jobCulturalFit, err
	}

	return jobCulturalFit, nil
}

func getJobRequirements(
	ctx context.Context,
	id int64,
	companyID int64,
) (model.JobRequirement, error) {
	jobRequirement := model.JobRequirement{}

	var itemsJSON []byte

	err := database.Database.QueryRowContext(
		ctx,
		`SELECT id,company_id,job_id,items,min_match,created_at,updated_at 
				FROM job_requirements WHERE company_id = ? AND job_id = ? LIMIT 1`,
		companyID,
		id,
	).Scan(
		&jobRequirement.ID,
		&jobRequirement.CompanyID,
		&jobRequirement.JobID,
		&itemsJSON,
		&jobRequirement.MinMatch,
		&jobRequirement.CreatedAt,
		&jobRequirement.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return jobRequirement, ErrJobNotFound
	}

	if err != nil {
		logrus.WithError(err).Error("Error to get job requirement")

		return jobRequirement, err
	}

	if err = json.Unmarshal(itemsJSON, &jobRequirement.Items); err != nil {
		logrus.WithError(err).Error("Error to unmarshal job requirement")

		return jobRequirement, err
	}

	return jobRequirement, nil
}

func getJobBenefits(
	ctx context.Context,
	id int64,
	companyID int64,
) ([]model.Benefit, error) {
	benefits := make([]model.Benefit, 0)

	rows, err := database.Database.QueryContext(
		ctx,
		`SELECT benefits.id, benefits.name, benefits.company_id, benefits.created_at, benefits.updated_at FROM job_benefits
				JOIN benefits ON job_benefits.benefit_id = benefits.id 
				WHERE job_benefits.company_id = ? AND job_benefits.job_id = ?`,
		companyID,
		id,
	)
	if err != nil {
		logrus.WithError(err).Error("Error to get job benefits")

		return benefits, err
	}

	defer rows.Close()

	for rows.Next() {
		benefit := model.Benefit{}

		err := rows.Scan(
			&benefit.ID,
			&benefit.Name,
			&benefit.CompanyID,
			&benefit.CreatedAt,
			&benefit.UpdatedAt,
		)
		if err != nil {
			logrus.WithError(err).Error("Error to scan job benefit")

			return benefits, err
		}

		benefits = append(benefits, benefit)
	}

	return benefits, nil
}
