package questionnaire_adp

import (
	"context"
	"encoding/json"

	"github.com/hubjob/api/database"
	"github.com/hubjob/api/model"
	"github.com/sirupsen/logrus"
)

func GetCandidateQuestionnaire(
	ctx context.Context,
	candidateID int64,
	questionnaireType string,
) (model.CandidateQuestionnaire, error) {
	questionnaire := model.CandidateQuestionnaire{}

	var answersString []byte

	err := database.Database.QueryRowContext(
		ctx,
		`SELECT id,candidate_id,type,answers,expired_at,created_at,updated_at 
				FROM candidate_questionnaires 
				WHERE candidate_id = ? AND type = ? 
				ORDER BY id DESC`,
		candidateID,
		questionnaireType,
	).Scan(
		&questionnaire.ID,
		&questionnaire.CandidateID,
		&questionnaire.Type,
		&answersString,
		&questionnaire.ExpiredAt,
		&questionnaire.CreatedAt,
		&questionnaire.UpdatedAt,
	)
	if err != nil {
		logrus.WithError(err).Error("Error to get candidate questionnaire")

		return questionnaire, err
	}

	err = json.Unmarshal(answersString, &questionnaire.Answers)
	if err != nil {
		logrus.WithError(err).Error("Error to unmarshal candidate questionnaire answers")

		return questionnaire, err
	}

	return questionnaire, nil
}
