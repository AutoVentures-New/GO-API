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

var ErrJobNotFound = errors.New("job not found")
var ErrInvalidJob = errors.New("invalid job")
var ErrApplicationAlreadyExist = errors.New("application already exist")

func StartApplication(
	ctx context.Context,
	candidateID,
	jobID int64,
) (model.Application, error) {
	job, err := getJob(ctx, jobID)
	if err != nil {
		return model.Application{}, err
	}

	if err := alreadyExist(ctx, candidateID, jobID); err != nil {
		return model.Application{}, err
	}

	hashQuestionnaire, err := hasCandidateQuestionnaire(ctx, candidateID, job.Questionnaire)
	if err != nil {
		return model.Application{}, err
	}

	questionnaire := model.QUESTIONNAIRE_BEHAVIORAL
	if job.Questionnaire == model.PROFESSIONAL {
		questionnaire = model.QUESTIONNAIRE_PROFESSIONAL
	}

	steps := []string{
		model.REQUIREMENTS,
		model.JOB_QUESTIONS,
		model.CULTURAL_FIT,
	}

	if job.Questionnaire != model.NONE && !hashQuestionnaire {
		steps = append(steps, questionnaire)
	}

	steps = append(steps, model.CANDIDATE_VIDEO)

	now := time.Now().UTC()
	application := model.Application{
		CompanyID:   job.CompanyID,
		JobID:       job.ID,
		CandidateID: candidateID,
		Steps:       steps,
		CurrentStep: model.REQUIREMENTS,
		Status:      model.FILLING,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	stepsString, err := json.Marshal(application.Steps)
	if err != nil {
		logrus.WithError(err).Error("Error to marshal application steps")

		return application, err
	}

	dbTransaction, err := database.Database.Begin()
	if err != nil {
		logrus.WithError(err).Error("Error to open db transaction")

		return application, err
	}

	result, err := dbTransaction.ExecContext(
		ctx,
		`INSERT INTO job_applications(company_id,job_id,candidate_id,steps,current_step,status,created_at,updated_at) 
					VALUES(?,?,?,?,?,?,?,?)`,
		application.CompanyID,
		application.JobID,
		application.CandidateID,
		stepsString,
		application.CurrentStep,
		application.Status,
		application.CreatedAt,
		application.UpdatedAt,
	)
	if err != nil {
		logrus.WithError(err).Error("Error to insert job application")

		_ = dbTransaction.Rollback()

		return application, err
	}

	application.ID, err = result.LastInsertId()
	if err != nil {
		logrus.WithError(err).Error("Error to get last insert job application id")

		_ = dbTransaction.Rollback()

		return application, err
	}

	if err := dbTransaction.Commit(); err != nil {
		_ = dbTransaction.Rollback()

		logrus.WithError(err).Error("Error to commit transaction")

		return application, err
	}

	return application, nil
}

func getJob(
	ctx context.Context,
	id int64,
) (model.Job, error) {
	job := model.Job{}

	err := database.Database.QueryRowContext(
		ctx,
		`SELECT id,company_id,status,questionnaire,publish_at,finish_at 
				FROM jobs WHERE id = ?`,
		id,
	).Scan(
		&job.ID,
		&job.CompanyID,
		&job.Status,
		&job.Questionnaire,
		&job.PublishAt,
		&job.FinishAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return job, ErrJobNotFound
	}

	if err != nil {
		logrus.WithError(err).Error("Error to get job")

		return job, err
	}

	now := time.Now().UTC()
	if now.Before(job.PublishAt) || now.After(job.FinishAt) {
		return job, ErrInvalidJob
	}

	return job, nil
}

func alreadyExist(
	ctx context.Context,
	candidateID,
	jobID int64,
) error {
	var count int

	err := database.Database.QueryRowContext(
		ctx,
		`SELECT COUNT(0) FROM job_applications WHERE candidate_id = ? AND job_id = ? LIMIT 1`,
		candidateID,
		jobID,
	).Scan(&count)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		logrus.WithError(err).Error("Error to validate job application")

		return err
	}

	if count > 0 {
		return ErrApplicationAlreadyExist
	}

	return nil
}

func hasCandidateQuestionnaire(
	ctx context.Context,
	candidateID int64,
	questionnaireType string,
) (bool, error) {
	var count int

	err := database.Database.QueryRowContext(
		ctx,
		`SELECT COUNT(0) FROM candidate_questionnaires WHERE candidate_id = ? AND type = ? AND expired_at > ? ORDER BY id DESC LIMIT 1`,
		candidateID,
		questionnaireType,
		time.Now().UTC(),
	).Scan(&count)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		logrus.WithError(err).Error("Error to validate candidate questionnaire")

		return false, err
	}

	if count > 0 {
		return true, nil
	}

	return false, nil
}
