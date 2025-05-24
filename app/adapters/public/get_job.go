package public

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
	candidateID int64,
) (model.Job, error) {
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
	if errors.Is(err, sql.ErrNoRows) {
		return job, ErrJobNotFound
	}

	if err != nil {
		logrus.WithError(err).Error("Error to get job")

		return job, err
	}

	job.Benefits, err = getJobBenefits(ctx, id)
	if err != nil {
		return job, err
	}

	job.Area, err = getArea(ctx, job.AreaID)
	if err != nil {
		return job, err
	}

	job.JobRequirement, err = getJobRequirements(ctx, id)
	if err != nil {
		return job, err
	}

	if candidateID > 0 {
		var count int64

		err := database.Database.QueryRowContext(
			ctx,
			`SELECT count(0) FROM job_applications WHERE candidate_id = ? AND job_id = ?`,
			candidateID,
			id,
		).Scan(&count)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			logrus.WithError(err).Error("Error to get job application for job details")

			return job, nil
		}

		if count > 0 {
			job.CandidateHasApplication = true
		}
	}

	return job, nil
}

func getJobBenefits(
	ctx context.Context,
	id int64,
) ([]model.Benefit, error) {
	benefits := make([]model.Benefit, 0)

	rows, err := database.Database.QueryContext(
		ctx,
		`SELECT benefits.id, benefits.name, benefits.company_id, benefits.created_at, benefits.updated_at FROM job_benefits
				JOIN benefits ON job_benefits.benefit_id = benefits.id 
				WHERE job_benefits.job_id = ?`,
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

func getArea(
	ctx context.Context,
	id int64,
) (*model.Area, error) {
	area := model.Area{}

	err := database.Database.QueryRowContext(
		ctx,
		`SELECT id,title,created_at,updated_at 
				FROM areas WHERE id = ?`,
		id,
	).Scan(
		&area.ID,
		&area.Title,
		&area.CreatedAt,
		&area.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrJobNotFound
	}

	if err != nil {
		logrus.WithError(err).Error("Error to get area")

		return nil, err
	}

	return &area, nil
}

func getJobRequirements(
	ctx context.Context,
	jobID int64,
) (*model.JobRequirement, error) {
	jobRequirement := model.JobRequirement{}

	var itemsJSON []byte

	err := database.Database.QueryRowContext(
		ctx,
		`SELECT items 
				FROM job_requirements WHERE job_id = ?`,
		jobID,
	).Scan(
		&itemsJSON,
	)
	if err != nil {
		logrus.WithError(err).Error("Error to get job requirement")

		return nil, err
	}

	if err = json.Unmarshal(itemsJSON, &jobRequirement.Items); err != nil {
		logrus.WithError(err).Error("Error to unmarshal job requirement")

		return nil, err
	}

	for index, _ := range jobRequirement.Items {
		jobRequirement.Items[index].Required = false
	}

	return &jobRequirement, nil
}
