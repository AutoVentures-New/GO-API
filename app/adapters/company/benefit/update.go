package benefit

import (
	"context"
	"time"

	"github.com/hubjob/api/database"
	"github.com/hubjob/api/model"
	"github.com/sirupsen/logrus"
)

func UpdateBenefit(
	ctx context.Context,
	benefit model.Benefit,
) (model.Benefit, error) {
	benefit.UpdatedAt = time.Now().UTC()

	_, err := database.Database.ExecContext(
		ctx,
		`UPDATE benefits set name = ?, updated_at = ? WHERE id = ?`,
		benefit.Name,
		benefit.UpdatedAt,
		benefit.ID,
	)
	if err != nil {
		logrus.WithError(err).Error("Error to update benefit")

		return benefit, err
	}

	return benefit, nil
}
