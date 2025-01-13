package questionnaire

import (
	"context"

	"github.com/hubjob/api/database"
	"github.com/hubjob/api/model"
	"github.com/sirupsen/logrus"
)

func ListQuestionnaires(
	ctx context.Context,
	companyID int64,
) ([]model.Questionnaire, error) {
	questionnaires := make([]model.Questionnaire, 0)

	rows, err := database.Database.QueryContext(
		ctx,
		`SELECT id, name, company_id, created_at, updated_at FROM questionnaires WHERE company_id = ?`,
		companyID,
	)
	if err != nil {
		logrus.WithError(err).Error("Error to list questionnaires")

		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		questionnaire := model.Questionnaire{}
		err := rows.Scan(
			&questionnaire.ID,
			&questionnaire.Name,
			&questionnaire.CompanyID,
			&questionnaire.CreatedAt,
			&questionnaire.UpdatedAt,
		)
		if err != nil {
			logrus.WithError(err).Error("Error to scan questionnaires")

			return nil, err
		}

		questionnaires = append(questionnaires, questionnaire)
	}

	return questionnaires, nil
}
