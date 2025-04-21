package queue_job

import (
	"context"

	"github.com/hubjob/api/app/adapters/queue_job"
	"github.com/hubjob/api/handler/queue_job/candidate_questionnaire"
	"github.com/hubjob/api/model"
	"github.com/sirupsen/logrus"
)

const (
	CANDIDATE_QUESTIONNAIRE_BEHAVIORAL   = "CANDIDATE-QUESTIONNAIRE-BEHAVIORAL"
	CANDIDATE_QUESTIONNAIRE_PROFESSIONAL = "CANDIDATE-QUESTIONNAIRE-PROFESSIONAL"
)

func Executor(ctx context.Context) {
	jobs, err := queue_job.ListJobs(ctx)
	if err != nil {
		logrus.WithError(err).Error("Error listing jobs")

		return
	}

	if len(jobs) == 0 {
		return
	}

	for _, job := range jobs {
		if job.Type == CANDIDATE_QUESTIONNAIRE_BEHAVIORAL {
			err = candidate_questionnaire.ExecuteQuestionnaireBehavioral(ctx, job)
		}

		if job.Type == CANDIDATE_QUESTIONNAIRE_PROFESSIONAL {
			err = candidate_questionnaire.ExecuteQuestionnaireProfessional(ctx, job)
		}

		job.Status = model.FINISHED
		if err != nil {
			job.Status = model.ERROR
		}

		err = queue_job.UpdateJob(ctx, job)
		if err != nil {
			logrus.WithError(err).Error("Error updating job")

			continue
		}
	}
}
