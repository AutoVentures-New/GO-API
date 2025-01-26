package steps

import (
	"context"
	"encoding/json"
	"time"

	"github.com/hubjob/api/database"
	"github.com/hubjob/api/model"
	"github.com/sirupsen/logrus"
)

func SaveJobQuestions(
	ctx context.Context,
	application model.Application,
) (model.Application, error) {
	dbTransaction, err := database.Database.Begin()
	if err != nil {
		logrus.WithError(err).Error("Error to open db transaction")

		return application, err
	}

	now := time.Now().UTC()
	jobApplicationQuestion := model.JobApplicationQuestion{
		ApplicationID: application.ID,
		Questions:     application.Questions,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	questionsString, err := json.Marshal(jobApplicationQuestion.Questions)
	if err != nil {
		logrus.WithError(err).Error("Error to marshal application questions")

		return application, err
	}

	_, err = dbTransaction.ExecContext(
		ctx,
		`INSERT INTO job_application_questions(application_id,questions,created_at,updated_at) VALUES(?,?,?,?)`,
		jobApplicationQuestion.ApplicationID,
		questionsString,
		jobApplicationQuestion.CreatedAt,
		jobApplicationQuestion.UpdatedAt,
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
