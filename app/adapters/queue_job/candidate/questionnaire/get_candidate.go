package questionnaire_adp

import (
	"context"

	"github.com/hubjob/api/database"
	"github.com/hubjob/api/model"
	"github.com/sirupsen/logrus"
)

func GetCandidate(
	ctx context.Context,
	candidateID int64,
) (model.Candidate, error) {
	var candidate model.Candidate

	err := database.Database.QueryRowContext(
		ctx,
		`SELECT id, name FROM candidates WHERE id = ?`,
		candidateID,
	).Scan(&candidate.ID, &candidate.Name)
	if err != nil {
		logrus.WithError(err).Error("Error to get candidate")

		return candidate, err
	}

	return candidate, nil
}
