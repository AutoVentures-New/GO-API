package job

import (
	"context"
	"database/sql"
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
		`UPDATE jobs set title = ?,area_id=?,is_talent_bank = ?,is_special_needs = ?,description = ?,job_mode = ?,
                contracting_modality = ?,state = ?,city = ?,responsibilities = ?,questionnaire = ?,video_link = ?,
                status = ?,publish_at = ?,updated_at = ? WHERE id = ?`,
		job.Title,
		job.AreaID,
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
		job.UpdatedAt,
		job.ID,
	)
	if err != nil {
		logrus.WithError(err).Error("Error to update job")

		_ = dbTransaction.Rollback()

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

	err = updateJobRequirement(ctx, dbTransaction, &job)
	if err != nil {
		_ = dbTransaction.Rollback()

		return job, err
	}

	if err := updateJobBenefits(ctx, dbTransaction, &job); err != nil {
		_ = dbTransaction.Rollback()

		return job, err
	}

	if err := updateJobVideoQuestions(ctx, dbTransaction, &job); err != nil {
		_ = dbTransaction.Rollback()

		return job, err
	}

	if err := updateJobQuestions(ctx, dbTransaction, &job); err != nil {
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

func updateJobRequirement(
	ctx context.Context,
	dbTransaction *sql.Tx,
	job *model.Job,
) error {
	job.JobRequirement.UpdatedAt = job.UpdatedAt

	itemsString, err := json.Marshal(job.JobRequirement.Items)
	if err != nil {
		logrus.WithError(err).Error("Error to marshal job requirements items")

		return err
	}

	_, err = dbTransaction.ExecContext(
		ctx,
		`UPDATE job_requirements set items = ?, min_match = ?, updated_at = ? WHERE id = ?`,
		itemsString,
		job.JobRequirement.MinMatch,
		job.JobRequirement.UpdatedAt,
		job.JobRequirement.ID,
	)
	if err != nil {
		logrus.WithError(err).Error("Error to update job requirements")

		return err
	}

	return nil
}

func updateJobBenefits(
	ctx context.Context,
	dbTransaction *sql.Tx,
	job *model.Job,
) error {
	_, err := dbTransaction.ExecContext(
		ctx,
		`DELETE FROM job_benefits WHERE company_id = ? AND job_id = ?`,
		job.CompanyID,
		job.ID,
	)
	if err != nil {
		logrus.WithError(err).Error("Error to delete job benefits")

		return err
	}

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

func updateJobVideoQuestions(
	ctx context.Context,
	dbTransaction *sql.Tx,
	job *model.Job,
) error {
	job.VideoQuestions.UpdatedAt = job.UpdatedAt

	questionsString, err := json.Marshal(job.VideoQuestions.Questions)
	if err != nil {
		logrus.WithError(err).Error("Error to marshal job video questions")

		return err
	}

	_, err = dbTransaction.ExecContext(
		ctx,
		`UPDATE job_video_questions SET questions = ?, updated_at = ? WHERE id = ?`,
		questionsString,
		job.VideoQuestions.UpdatedAt,
		job.VideoQuestions.ID,
	)
	if err != nil {
		logrus.WithError(err).Error("Error to update job video questions")

		return err
	}

	return nil
}

func updateJobQuestions(
	ctx context.Context,
	dbTransaction *sql.Tx,
	job *model.Job,
) error {
	_, err := dbTransaction.ExecContext(
		ctx,
		`DELETE FROM job_questions WHERE company_id = ? AND job_id = ?`,
		job.CompanyID,
		job.ID,
	)
	if err != nil {
		logrus.WithError(err).Error("Error to delete job questions")

		return err
	}

	for _, question := range job.Questions {
		_, err := dbTransaction.ExecContext(
			ctx,
			`INSERT INTO job_questions(company_id,job_id,question_id,created_at,updated_at) 
					VALUES(?,?,?,?,?)`,
			job.CompanyID,
			job.ID,
			question.ID,
			job.CreatedAt,
			job.UpdatedAt,
		)
		if err != nil {
			logrus.WithError(err).Error("Error to insert job questions")

			return err
		}
	}

	return nil
}
