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

var ErrApplicationNotFound = errors.New("application not found")

func GetJobApplication(
	ctx context.Context,
	jobID int64,
	candidateID int64,
) (model.Application, error) {
	application := model.Application{}

	var stepsJSON []byte

	err := database.Database.QueryRowContext(
		ctx,
		`SELECT id,company_id,job_id,candidate_id,steps,current_step,status,created_at,updated_at 
				FROM job_applications WHERE candidate_id = ? AND job_id = ?`,
		candidateID,
		jobID,
	).Scan(
		&application.ID,
		&application.CompanyID,
		&application.JobID,
		&application.CandidateID,
		&stepsJSON,
		&application.CurrentStep,
		&application.Status,
		&application.CreatedAt,
		&application.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return application, ErrApplicationNotFound
	}

	if err != nil {
		logrus.WithError(err).Error("Error to get job application")

		return application, err
	}

	err = json.Unmarshal(stepsJSON, &application.Steps)
	if err != nil {
		logrus.WithError(err).Error("Error to unmarshal application steps")

		return application, err
	}

	if application.Status != model.FILLING {
		return application, nil
	}

	if application.CurrentStep == model.REQUIREMENTS {
		application.JobApplicationRequirementItem, err = getJobRequirements(ctx, application.JobID, application.CompanyID)
		if err != nil {
			return application, err
		}
	}

	return application, nil
}

func getJobRequirements(
	ctx context.Context,
	jobID int64,
	companyID int64,
) ([]model.JobApplicationRequirementItem, error) {
	appRequirements := make([]model.JobApplicationRequirementItem, 0)
	jobRequirement := model.JobRequirement{}

	var itemsJSON []byte

	err := database.Database.QueryRowContext(
		ctx,
		`SELECT items 
				FROM job_requirements WHERE company_id = ? AND job_id = ?`,
		companyID,
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

	for _, jobRequirementItem := range jobRequirement.Items {
		appRequirements = append(appRequirements, model.JobApplicationRequirementItem{
			ID:      jobRequirementItem.ID,
			Name:    jobRequirementItem.Name,
			Checked: false,
		})
	}

	return appRequirements, nil
}
