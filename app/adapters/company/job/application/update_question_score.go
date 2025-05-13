package job

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"math"
	"time"

	"github.com/hubjob/api/database"
	"github.com/hubjob/api/model"
	"github.com/sirupsen/logrus"
)

func UpdateQuestionScore(
	ctx context.Context,
	applicationID int64,
	questionID int64,
	score int64,
) error {
	var jobApplicationQuestion model.JobApplicationQuestion

	var questionsString []byte

	err := database.Database.QueryRowContext(
		ctx,
		`SELECT application_id,questions,score,open_field_score,created_at,updated_at FROM job_application_questions WHERE application_id = ?`,
		applicationID,
	).Scan(
		&jobApplicationQuestion.ApplicationID,
		&questionsString,
		&jobApplicationQuestion.Score,
		&jobApplicationQuestion.OpenFieldScore,
		&jobApplicationQuestion.CreatedAt,
		&jobApplicationQuestion.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil
	}

	if err != nil {
		logrus.WithError(err).Error("Error to get job application question")

		return err
	}

	err = json.Unmarshal(questionsString, &jobApplicationQuestion.Questions)
	if err != nil {
		logrus.WithError(err).Error("Error to unmarshal application question")

		return err
	}

	var total, count int64

	for index, question := range jobApplicationQuestion.Questions {
		if question.ID != questionID || question.Type != model.OPEN_FIELD {
			continue
		}

		jobApplicationQuestion.Questions[index].Score = score
	}

	for _, question := range jobApplicationQuestion.Questions {
		if question.Type != model.OPEN_FIELD {
			continue
		}

		total += question.Score
		if question.Score > 0 {
			count++
		}
	}

	result := total

	if count > 0 {
		result = int64(math.Round(float64(result) / float64(count)))
	}

	jobApplicationQuestion.OpenFieldScore = result

	questionsJson, err := json.Marshal(jobApplicationQuestion.Questions)
	if err != nil {
		logrus.WithError(err).Error("Error to marshal job application question")

		return err
	}

	_, err = database.Database.ExecContext(
		ctx,
		`UPDATE job_application_questions set questions = ?, open_field_score = ?, updated_at = ? WHERE application_id = ?`,
		questionsJson,
		jobApplicationQuestion.OpenFieldScore,
		time.Now().UTC(),
		applicationID,
	)
	if err != nil {
		logrus.WithError(err).Error("Error to update job application question")

		return err
	}

	return nil
}
