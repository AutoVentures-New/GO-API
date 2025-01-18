package job

import (
	"context"
	"encoding/json"
	"time"

	"github.com/hubjob/api/database"
	"github.com/hubjob/api/model"
	"github.com/sirupsen/logrus"
)

func UpdateJob(
	ctx context.Context,
	job model.Job,
) (model.Job, error) {
	dbTransaction, err := database.Database.Begin()
	if err != nil {
		logrus.WithError(err).Error("Error to open db transaction")

		return job, err
	}

	job.UpdatedAt = time.Now().UTC()

	_, err = dbTransaction.ExecContext(
		ctx,
		`UPDATE jobs set title = ?,is_talent_bank = ?,is_special_needs = ?,description = ?,job_mode = ?,
                contracting_modality = ?,state = ?,city = ?,responsibilities = ?,questionnaire = ?,video_link = ?,
                status = ?,publish_at = ?,finish_at = ?,updated_at = ? WHERE id = ?`,
		job.Title,
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
		job.UpdatedAt,
		job.ID,
	)
	if err != nil {
		logrus.WithError(err).Error("Error to update job")

		return job, err
	}

	job.JobCulturalFit.UpdatedAt = job.UpdatedAt

	answersString, err := json.Marshal(job.JobCulturalFit.Answers)
	if err != nil {
		logrus.WithError(err).Error("Error to marshal job cultural fit answers")

		_ = dbTransaction.Rollback()

		return job, err
	}

	_, err = dbTransaction.ExecContext(
		ctx,
		`UPDATE job_cultural_fit set answers = ?,updated_at = ? WHERE id = ?`,
		answersString,
		job.JobCulturalFit.UpdatedAt,
		job.JobCulturalFit.ID,
	)
	if err != nil {
		logrus.WithError(err).Error("Error to insert job cultural fit")

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
