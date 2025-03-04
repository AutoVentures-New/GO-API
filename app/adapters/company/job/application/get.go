package job

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	profile "github.com/hubjob/api/app/adapters/candidate/curriculum"
	"github.com/hubjob/api/database"
	"github.com/hubjob/api/model"
	"github.com/sirupsen/logrus"
)

var ErrApplicationNotFound = errors.New("Application not found")

func GetApplication(
	ctx context.Context,
	companyID int64,
	id int64,
) (model.Application, error) {
	application := model.Application{}

	var stepsJSON []byte

	err := database.Database.QueryRowContext(
		ctx,
		`SELECT id,company_id,job_id,candidate_id,steps,current_step,status,created_at,updated_at 
				FROM job_applications WHERE company_id = ? AND id = ?`,
		companyID,
		id,
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

	application.Candidate, err = getCandidate(ctx, application.CandidateID)
	if err != nil {
		return application, err
	}

	application.JobApplicationRequirement, err = getJobApplicationRequirement(ctx, application.ID)
	if err != nil {
		return application, err
	}

	return application, nil
}

func getCandidate(
	ctx context.Context,
	candidateID int64,
) (*model.Candidate, error) {
	var candidate model.Candidate

	err := database.Database.QueryRowContext(
		ctx,
		`SELECT id,name,cpf,email,password,status,phone,birth_date,created_at,updated_at FROM candidates WHERE id = ?`,
		candidateID,
	).Scan(
		&candidate.ID,
		&candidate.Name,
		&candidate.CPF,
		&candidate.Email,
		&candidate.Password,
		&candidate.Status,
		&candidate.Phone,
		&candidate.BirthDate,
		&candidate.CreatedAt,
		&candidate.UpdatedAt,
	)
	if err != nil {
		logrus.WithError(err).Error("Error to get candidate")

		return nil, err
	}

	_ = database.Database.QueryRowContext(
		ctx,
		`SELECT state, city FROM candidate_addresses WHERE candidate_id = ?`,
		candidateID,
	).Scan(
		&candidate.Address.State,
		&candidate.Address.City,
	)

	curriculum, err := profile.GetCurriculum(ctx, candidate.ID)
	if err != nil {
		return nil, err
	}

	candidate.Curriculum = &curriculum

	return &candidate, nil
}

func getJobApplicationRequirement(
	ctx context.Context,
	applicationID int64,
) (*model.JobApplicationRequirement, error) {
	var jobApplicationRequirement model.JobApplicationRequirement

	var itemsString []byte

	err := database.Database.QueryRowContext(
		ctx,
		`SELECT application_id,items,match_value,created_at,updated_at FROM job_application_requirements WHERE application_id = ?`,
		applicationID,
	).Scan(
		&jobApplicationRequirement.ApplicationID,
		&itemsString,
		&jobApplicationRequirement.MatchValue,
		&jobApplicationRequirement.CreatedAt,
		&jobApplicationRequirement.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		logrus.WithError(err).Error("Error to get job application requirements")

		return nil, err
	}

	err = json.Unmarshal(itemsString, &jobApplicationRequirement.Items)
	if err != nil {
		logrus.WithError(err).Error("Error to unmarshal application requirements")

		return nil, err
	}

	return &jobApplicationRequirement, nil
}
