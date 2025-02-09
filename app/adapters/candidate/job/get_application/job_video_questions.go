package get_application

import (
	"context"
	"encoding/json"

	"github.com/hubjob/api/database"
	"github.com/hubjob/api/model"
	"github.com/sirupsen/logrus"
)

func GetJobVideoQuestions(
	ctx context.Context,
	jobID int64,
	companyID int64,
) (*model.JobVideoQuestions, error) {
	jobVideoQuestions := model.JobVideoQuestions{}

	var itemsJSON []byte

	err := database.Database.QueryRowContext(
		ctx,
		`SELECT questions 
				FROM job_video_questions WHERE company_id = ? AND job_id = ?`,
		companyID,
		jobID,
	).Scan(&itemsJSON)
	if err != nil {
		logrus.WithError(err).Error("Error to get job video questions")

		return nil, err
	}

	if err = json.Unmarshal(itemsJSON, &jobVideoQuestions.Questions); err != nil {
		logrus.WithError(err).Error("Error to unmarshal job video questions")

		return nil, err
	}

	return &jobVideoQuestions, nil
}
