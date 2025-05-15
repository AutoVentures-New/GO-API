package steps

import (
	"context"
	"encoding/json"
	"math"
	"time"

	questions_adp "github.com/hubjob/api/app/adapters/questions"
	"github.com/hubjob/api/database"
	"github.com/hubjob/api/model"
	"github.com/sirupsen/logrus"
)

func SaveJobQuestions(
	ctx context.Context,
	application model.Application,
) (model.Application, error) {
	questions, err := questions_adp.ListQuestions(ctx, application.JobID)
	if err != nil {
		return application, err
	}

	var score, correct int64

	for _, questionC := range application.Questions {
		question, ok := questions[questionC.ID]
		if !ok {
			continue
		}

		if question.Type == model.OPEN_FIELD {
			continue
		}

		if question.Type == model.SINGLE_CHOICE {
			var correctAnswer int64

			for _, answer := range question.Answers {
				if answer.IsCorrect {
					correctAnswer = answer.ID
					correct++
				}
			}

			for _, answer := range questionC.Answers {
				if !answer.Checked {
					continue
				}

				if answer.ID == correctAnswer {
					score++
				}
			}

			continue
		}

		correctAnswer := make(map[int64]int64)

		for _, answer := range question.Answers {
			if answer.IsCorrect {
				correctAnswer[answer.ID] = answer.ID
				correct++
			}
		}

		for _, answer := range questionC.Answers {
			if !answer.Checked {
				continue
			}

			if _, ok := correctAnswer[answer.ID]; !ok {
				continue
			}

			score++
		}
	}

	dbTransaction, err := database.Database.Begin()
	if err != nil {
		logrus.WithError(err).Error("Error to open db transaction")

		return application, err
	}

	percent := (float64(score) / float64(correct)) * 100

	now := time.Now().UTC()
	jobApplicationQuestion := model.JobApplicationQuestion{
		ApplicationID: application.ID,
		Questions:     application.Questions,
		Score:         int64(math.Round(percent)),
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	questionsString, err := json.Marshal(jobApplicationQuestion.Questions)
	if err != nil {
		logrus.WithError(err).Error("Error to marshal application questions")

		return application, err
	}

	_, err = dbTransaction.ExecContext(
		ctx,
		`INSERT INTO job_application_questions(application_id,questions,score,open_field_score,created_at,updated_at) VALUES(?,?,?,?,?,?)`,
		jobApplicationQuestion.ApplicationID,
		questionsString,
		jobApplicationQuestion.Score,
		jobApplicationQuestion.OpenFieldScore,
		jobApplicationQuestion.CreatedAt,
		jobApplicationQuestion.UpdatedAt,
	)
	if err != nil {
		logrus.WithError(err).Error("Error to insert job application questions")

		_ = dbTransaction.Rollback()

		return application, err
	}

	if err := updateApplication(ctx, dbTransaction, &application, false); err != nil {
		_ = dbTransaction.Rollback()

		return application, err
	}

	if err := dbTransaction.Commit(); err != nil {
		_ = dbTransaction.Rollback()

		logrus.WithError(err).Error("Error to commit transaction")

		return application, err
	}

	return application, nil
}
