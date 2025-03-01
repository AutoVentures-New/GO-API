package profile

import (
	"context"

	"github.com/hubjob/api/database"
	"github.com/hubjob/api/model"
	"github.com/sirupsen/logrus"
)

func GetCompany(
	ctx context.Context,
	companyID int64,
) (model.Company, error) {
	company := model.Company{}

	var desc, logo *string

	err := database.Database.QueryRowContext(
		ctx,
		`SELECT id,name,cnpj,description,logo,status,created_at,updated_at 
				FROM companies WHERE id = ?`,
		companyID,
	).Scan(
		&company.ID,
		&company.Name,
		&company.CNPJ,
		&desc,
		&logo,
		&company.Status,
		&company.CreatedAt,
		&company.UpdatedAt,
	)
	if err != nil {
		logrus.WithError(err).Error("Error to get company")

		return company, err
	}

	if desc != nil {
		company.Description = *desc
	}

	if logo != nil {
		company.Logo = *logo
	}

	return company, nil
}
