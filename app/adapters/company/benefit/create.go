package benefit

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/hubjob/api/database"
	"github.com/hubjob/api/model"
	"github.com/sirupsen/logrus"
)

var ErrBenefitAlreadyExists = errors.New("benefit already exists")

func CreateBenefit(
	ctx context.Context,
	benefit model.Benefit,
) (model.Benefit, error) {
	var count int

	err := database.Database.QueryRowContext(
		ctx,
		`SELECT COUNT(0) FROM benefits WHERE company_id = ? AND name = ?`,
		benefit.CompanyID,
		benefit.Name,
	).Scan(&count)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		logrus.WithError(err).Error("Error to validate benefit")

		return benefit, err
	}

	if count > 0 {
		return benefit, ErrBenefitAlreadyExists
	}

	benefit.CreatedAt = time.Now().UTC()
	benefit.UpdatedAt = benefit.CreatedAt

	result, err := database.Database.ExecContext(
		ctx,
		`INSERT INTO benefits(name,company_id,created_at,updated_at) VALUES(?,?,?,?)`,
		benefit.Name,
		benefit.CompanyID,
		benefit.CreatedAt,
		benefit.UpdatedAt,
	)
	if err != nil {
		logrus.WithError(err).Error("Error to insert benefit")

		return benefit, err
	}

	benefit.ID, err = result.LastInsertId()
	if err != nil {
		logrus.WithError(err).Error("Error to get last insert benefit id")

		return benefit, err
	}

	return benefit, nil
}
