package job

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/hubjob/api/database"
	"github.com/hubjob/api/model"
	"github.com/sirupsen/logrus"
)

func GetLastCulturalFit(
	ctx context.Context,
	companyID int64,
) (model.CompanyCulturalFit, error) {
	companyCulturalFit := model.CompanyCulturalFit{}

	var answersJSON []byte

	err := database.Database.QueryRowContext(
		ctx,
		`SELECT company_id,answers 
				FROM job_cultural_fit WHERE company_id = ?
				ORDER BY id DESC LIMIT 1`,
		companyID,
	).Scan(
		&companyCulturalFit.CompanyID,
		&answersJSON,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return companyCulturalFit, nil
	}

	if err != nil {
		logrus.WithError(err).Error("Error to get job cultural fit")

		return companyCulturalFit, err
	}

	err = json.Unmarshal(answersJSON, &companyCulturalFit.Answers)
	if err != nil {
		logrus.WithError(err).Error("Error to unmarshal job cultural fit")

		return companyCulturalFit, err
	}

	return companyCulturalFit, nil
}
