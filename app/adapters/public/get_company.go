package public

import (
	"context"
	"database/sql"
	"errors"

	"github.com/hubjob/api/database"
	"github.com/hubjob/api/model"
	"github.com/sirupsen/logrus"
)

var ErrNotFound = errors.New("Company not found")

func GetCompany(
	ctx context.Context,
	companyID int64,
) (model.Company, error) {
	company := model.Company{}

	err := database.Database.QueryRowContext(
		ctx,
		`SELECT id,name,cnpj,description,logo,status,created_at,updated_at 
				FROM companies WHERE id = ?`,
		companyID,
	).Scan(
		&company.ID,
		&company.Name,
		&company.CNPJ,
		&company.Description,
		&company.Logo,
		&company.Status,
		&company.CreatedAt,
		&company.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return company, ErrNotFound
	}

	if err != nil {
		logrus.WithError(err).Error("Error to get company public")

		return company, err
	}

	return company, nil
}
