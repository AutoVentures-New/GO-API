package benefit

import (
	"context"
	"database/sql"
	"errors"

	"github.com/hubjob/api/database"
	"github.com/hubjob/api/model"
	"github.com/sirupsen/logrus"
)

var ErrNotFound = errors.New("Benefit not found")

func GetBenefit(
	ctx context.Context,
	id int64,
	companyID int64,
) (model.Benefit, error) {
	benefit := model.Benefit{}

	err := database.Database.QueryRowContext(
		ctx,
		`SELECT id, name, company_id, created_at, updated_at FROM benefits WHERE company_id = ? AND id = ?`,
		companyID,
		id,
	).Scan(
		&benefit.ID,
		&benefit.Name,
		&benefit.CompanyID,
		&benefit.CreatedAt,
		&benefit.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return benefit, ErrNotFound
	}

	if err != nil {
		logrus.WithError(err).Error("Error to get benefit")

		return benefit, err
	}

	return benefit, nil
}
