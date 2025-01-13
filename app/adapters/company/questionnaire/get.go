package questionnaire

import (
	"context"
	"database/sql"
	"errors"

	"github.com/hubjob/api/database"
	"github.com/hubjob/api/model"
	"github.com/sirupsen/logrus"
)

var ErrQuestionnaireNotFound = errors.New("Questionnaire not found")

func GetQuestionnaire(
	ctx context.Context,
	id int64,
	companyID int64,
) (model.Questionnaire, error) {
	questionnaire := model.Questionnaire{}

	err := database.Database.QueryRowContext(
		ctx,
		`SELECT id, name, company_id, created_at, updated_at FROM questionnaires WHERE company_id = ? AND id = ?`,
		companyID,
		id,
	).Scan(
		&questionnaire.ID,
		&questionnaire.Name,
		&questionnaire.CompanyID,
		&questionnaire.CreatedAt,
		&questionnaire.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return questionnaire, ErrQuestionnaireNotFound
	}

	if err != nil {
		logrus.WithError(err).Error("Error to get questionnaire")

		return questionnaire, err
	}

	return questionnaire, nil
}
