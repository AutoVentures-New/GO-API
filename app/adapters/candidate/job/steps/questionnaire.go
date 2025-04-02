package steps

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hubjob/api/database"
	"github.com/hubjob/api/model"
	"github.com/sirupsen/logrus"
)

func SaveQuestionnaire(
	ctx context.Context,
	application model.Application,
	questionnaire model.CandidateQuestionnaire,
) (model.Application, error) {
	dbTransaction, err := database.Database.Begin()
	if err != nil {
		logrus.WithError(err).Error("Error to open db transaction")

		return application, err
	}

	now := time.Now().UTC()
	questionnaire.CandidateID = application.CandidateID
	questionnaire.CreatedAt = now
	questionnaire.UpdatedAt = now
	questionnaire.ExpiredAt = now.AddDate(0, 6, 0)

	questionnaire.Type = model.BEHAVIORAL
	if application.CurrentStep == model.QUESTIONNAIRE_PROFESSIONAL {
		questionnaire.Type = model.PROFESSIONAL
	}

	answersString, err := json.Marshal(questionnaire.Answers)
	if err != nil {
		logrus.WithError(err).Error("Error to marshal candidate answers")

		return application, err
	}

	_, err = dbTransaction.ExecContext(
		ctx,
		`INSERT INTO candidate_questionnaires(candidate_id,type,answers,expired_at,created_at,updated_at) VALUES(?,?,?,?,?,?)`,
		questionnaire.CandidateID,
		questionnaire.Type,
		answersString,
		questionnaire.ExpiredAt,
		questionnaire.CreatedAt,
		questionnaire.UpdatedAt,
	)
	if err != nil {
		logrus.WithError(err).Error("Error to insert candidate questionnaires")

		_ = dbTransaction.Rollback()

		return application, err
	}

	queueJob := model.QueueJob{
		Type:   fmt.Sprintf("CANDIDATE-QUESTIONNAIRE-%s", questionnaire.Type),
		Status: model.PENDING_JOB,
		Configurations: model.Configurations{
			CandidateID: questionnaire.CandidateID,
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	configurationsString, err := json.Marshal(queueJob.Configurations)
	if err != nil {
		logrus.WithError(err).Error("Error to marshal queue job configurations")

		return application, err
	}

	_, err = dbTransaction.ExecContext(
		ctx,
		`INSERT INTO queue_jobs(type,status,configurations,created_at,updated_at) VALUES(?,?,?,?,?)`,
		queueJob.Type,
		queueJob.Status,
		configurationsString,
		queueJob.CreatedAt,
		queueJob.UpdatedAt,
	)
	if err != nil {
		logrus.WithError(err).Error("Error to insert queue job")

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
