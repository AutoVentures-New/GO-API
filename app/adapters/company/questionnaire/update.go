package questionnaire

import (
	"context"
	"time"

	"github.com/hubjob/api/database"
	"github.com/hubjob/api/model"
	"github.com/sirupsen/logrus"
)

func UpdateQuestionnaire(
	ctx context.Context,
	questionnaire model.Questionnaire,
) (model.Questionnaire, error) {
	questionnaire.UpdatedAt = time.Now().UTC()

	_, err := database.Database.ExecContext(
		ctx,
		`UPDATE questionnaires set name = ?, updated_at = ? WHERE id = ?`,
		questionnaire.Name,
		questionnaire.UpdatedAt,
		questionnaire.ID,
	)
	if err != nil {
		logrus.WithError(err).Error("Error to update questionnaire")

		return questionnaire, err
	}

	return questionnaire, nil
}
