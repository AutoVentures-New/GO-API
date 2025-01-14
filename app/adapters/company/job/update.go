package job

import (
	"context"
	"time"

	"github.com/hubjob/api/database"
	"github.com/hubjob/api/model"
	"github.com/sirupsen/logrus"
)

func UpdateJob(
	ctx context.Context,
	job model.Job,
) (model.Job, error) {
	job.UpdatedAt = time.Now().UTC()

	_, err := database.Database.ExecContext(
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

	return job, nil
}
