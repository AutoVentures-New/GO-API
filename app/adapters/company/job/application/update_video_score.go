package job

import (
	"context"
	"time"

	"github.com/hubjob/api/database"
	"github.com/sirupsen/logrus"
)

func UpdateCandidateVideoScore(
	ctx context.Context,
	applicationID int64,
	score int64,
) error {
	_, err := database.Database.ExecContext(
		ctx,
		`UPDATE job_application_candidate_videos set score = ?, updated_at = ? WHERE application_id = ?`,
		score,
		time.Now().UTC(),
		applicationID,
	)
	if err != nil {
		logrus.WithError(err).Error("Error to update candidate video score")

		return err
	}

	return nil
}
