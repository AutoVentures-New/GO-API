package candidate_questionnaire

import (
	"context"

	"github.com/hubjob/api/app/adapters/queue_job"
	"github.com/hubjob/api/model"
	"github.com/sirupsen/logrus"
)

func ExecuteQuestionnaireProfessional(
	ctx context.Context,
	queueJob model.QueueJob,
) {
	logrus.WithFields(logrus.Fields{
		"queue_job": queueJob,
	}).Info("Executing questionnaire professional")

	queueJob.Status = model.FINISHED

	err := queue_job.UpdateJob(ctx, queueJob)
	if err != nil {
		logrus.WithError(err).Error("Error updating job")

		return
	}
}
