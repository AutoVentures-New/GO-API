package questionnaire

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/hubjob/api/database"
	"github.com/hubjob/api/model"
	"github.com/sirupsen/logrus"
)

var ErrQuestionnaireAlreadyExists = errors.New("questionnaire already exists")

func CreateQuestionnaire(
	ctx context.Context,
	questionnaire model.Questionnaire,
) (model.Questionnaire, error) {
	var count int

	err := database.Database.QueryRowContext(
		ctx,
		`SELECT COUNT(0) FROM questionnaires WHERE company_id = ? AND name = ?`,
		questionnaire.CompanyID,
		questionnaire.Name,
	).Scan(&count)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		logrus.WithError(err).Error("Error to validate questionnaire")

		return questionnaire, err
	}

	if count > 0 {
		return questionnaire, ErrQuestionnaireAlreadyExists
	}

	questionnaire.CreatedAt = time.Now().UTC()
	questionnaire.UpdatedAt = questionnaire.CreatedAt

	result, err := database.Database.ExecContext(
		ctx,
		`INSERT INTO questionnaires(name,company_id,created_at,updated_at) VALUES(?,?,?,?)`,
		questionnaire.Name,
		questionnaire.CompanyID,
		questionnaire.CreatedAt,
		questionnaire.UpdatedAt,
	)
	if err != nil {
		logrus.WithError(err).Error("Error to insert questionnaire")

		return questionnaire, err
	}

	questionnaire.ID, err = result.LastInsertId()
	if err != nil {
		logrus.WithError(err).Error("Error to get last insert questionnaire id")

		return questionnaire, err
	}

	return questionnaire, nil
}
