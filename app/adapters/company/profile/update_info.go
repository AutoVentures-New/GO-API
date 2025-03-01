package profile

import (
	"context"
	"time"

	"github.com/hubjob/api/database"
	"github.com/hubjob/api/model"
	"github.com/sirupsen/logrus"
)

func UpdateInfo(
	ctx context.Context,
	company model.Company,
) (model.Company, error) {
	dbTransaction, err := database.Database.Begin()
	if err != nil {
		logrus.WithError(err).Error("Error to open db transaction")

		return model.Company{}, err
	}

	company.UpdatedAt = time.Now().UTC()

	_, err = dbTransaction.ExecContext(
		ctx,
		`UPDATE companies set
    				name = ?, 
    				cnpj = ?,  
    				description = ?,  
    				updated_at = ?
				WHERE id = ?`,
		company.Name,
		company.CNPJ,
		company.Description,
		company.UpdatedAt,
		company.ID,
	)
	if err != nil {
		_ = dbTransaction.Rollback()

		logrus.WithError(err).Error("Error to update company")

		return model.Company{}, err
	}

	err = dbTransaction.Commit()
	if err != nil {
		_ = dbTransaction.Rollback()

		logrus.WithError(err).Error("Error to commit db transaction")

		return model.Company{}, err
	}

	return company, nil
}
