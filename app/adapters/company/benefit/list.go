package benefit

import (
	"context"

	"github.com/hubjob/api/database"
	"github.com/hubjob/api/model"
	"github.com/sirupsen/logrus"
)

func ListBenefits(
	ctx context.Context,
	companyID int64,
) ([]model.Benefit, error) {
	benefits := make([]model.Benefit, 0)

	rows, err := database.Database.QueryContext(
		ctx,
		`SELECT id, name, company_id, created_at, updated_at FROM benefits WHERE company_id = ?`,
		companyID,
	)
	if err != nil {
		logrus.WithError(err).Error("Error to list benefits")

		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		benefit := model.Benefit{}
		err := rows.Scan(
			&benefit.ID,
			&benefit.Name,
			&benefit.CompanyID,
			&benefit.CreatedAt,
			&benefit.UpdatedAt,
		)
		if err != nil {
			logrus.WithError(err).Error("Error to scan benefits")

			return nil, err
		}

		benefits = append(benefits, benefit)
	}

	return benefits, nil
}
