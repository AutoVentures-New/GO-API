package public

import (
	"context"

	"github.com/hubjob/api/database"
	"github.com/hubjob/api/model"
	"github.com/sirupsen/logrus"
)

type ListCompaniesFilter struct {
}

func ListCompanies(
	ctx context.Context,
	filter ListCompaniesFilter,
) ([]model.Company, int64, error) {
	companies := make([]model.Company, 0)

	queryCount := `SELECT count(*) FROM companies WHERE status = 'ACTIVE'`

	query := `SELECT id,name,description,logo,created_at,updated_at 
				FROM companies WHERE status = 'ACTIVE'
				ORDER BY name ASC`

	var count int64

	err := database.Database.QueryRowContext(
		ctx,
		queryCount,
	).Scan(
		&count,
	)
	if err != nil {
		logrus.WithError(err).Error("Error to get companies count")

		return nil, 0, err
	}

	rows, err := database.Database.QueryContext(
		ctx,
		query,
	)
	if err != nil {
		logrus.WithError(err).Error("Error to list companies")

		return nil, 0, err
	}

	defer rows.Close()

	for rows.Next() {
		company := model.Company{}
		err := rows.Scan(
			&company.ID,
			&company.Name,
			&company.Description,
			&company.Logo,
			&company.CreatedAt,
			&company.UpdatedAt,
		)
		if err != nil {
			logrus.WithError(err).Error("Error to scan company")

			return nil, 0, err
		}

		companies = append(companies, company)
	}

	return companies, count, nil
}
