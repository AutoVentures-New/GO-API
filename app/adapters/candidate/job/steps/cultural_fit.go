package steps

import (
	"context"
	"encoding/json"
	"time"

	"github.com/hubjob/api/database"
	"github.com/hubjob/api/model"
	"github.com/sirupsen/logrus"
)

func SaveCulturalFit(
	ctx context.Context,
	application model.Application,
) (model.Application, error) {
	jobCulturalFit, err := getJobCulturalFit(ctx, application.JobID, application.CompanyID)
	if err != nil {
		return application, err
	}

	match := matchCulturalFitValue(application, jobCulturalFit)

	dbTransaction, err := database.Database.Begin()
	if err != nil {
		logrus.WithError(err).Error("Error to open db transaction")

		return application, err
	}

	now := time.Now().UTC()
	jobApplicationCulturalFit := model.JobApplicationCulturalFit{
		ApplicationID: application.ID,
		Answers:       application.CulturalFit.Answers,
		MatchValue:    match,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	answersString, err := json.Marshal(jobApplicationCulturalFit.Answers)
	if err != nil {
		logrus.WithError(err).Error("Error to marshal application cultural answers")

		return application, err
	}

	_, err = dbTransaction.ExecContext(
		ctx,
		`INSERT INTO job_application_cultural_fit(application_id,answers,match_value,created_at,updated_at) VALUES(?,?,?,?,?)`,
		jobApplicationCulturalFit.ApplicationID,
		answersString,
		jobApplicationCulturalFit.MatchValue,
		jobApplicationCulturalFit.CreatedAt,
		jobApplicationCulturalFit.UpdatedAt,
	)
	if err != nil {
		logrus.WithError(err).Error("Error to insert job application questions")

		_ = dbTransaction.Rollback()

		return application, err
	}

	if err := updateApplication(ctx, dbTransaction, &application, false); err != nil {
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

func matchCulturalFitValue(
	application model.Application,
	jobCulturalFit model.JobCulturalFit,
) int64 {
	culturalFitAnswers := make(map[int64]string)
	match := 0

	for _, answer := range jobCulturalFit.Answers {
		culturalFitAnswers[answer.CulturalFitID] = answer.Answer
	}

	for _, answer := range application.CulturalFit.Answers {
		value, ok := culturalFitAnswers[answer.CulturalFitID]
		if !ok {
			continue
		}

		if value != answer.Answer {
			continue
		}

		match++
	}

	return int64((match * 100) / len(jobCulturalFit.Answers))
}

func getJobCulturalFit(
	ctx context.Context,
	jobID int64,
	companyID int64,
) (model.JobCulturalFit, error) {
	jobCulturalFit := model.JobCulturalFit{}

	var answersJSON []byte

	err := database.Database.QueryRowContext(
		ctx,
		`SELECT answers 
				FROM job_cultural_fit WHERE company_id = ? AND job_id = ?`,
		companyID,
		jobID,
	).Scan(&answersJSON)
	if err != nil {
		logrus.WithError(err).Error("Error to get job cultural answers")

		return jobCulturalFit, err
	}

	if err = json.Unmarshal(answersJSON, &jobCulturalFit.Answers); err != nil {
		logrus.WithError(err).Error("Error to unmarshal job cultural answers")

		return jobCulturalFit, err
	}

	return jobCulturalFit, nil
}
