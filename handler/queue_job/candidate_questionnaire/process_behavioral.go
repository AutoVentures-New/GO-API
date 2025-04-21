package candidate_questionnaire

import (
	"context"

	"github.com/hubjob/api/model"
	"github.com/sirupsen/logrus"
)

func ExecuteQuestionnaireBehavioral(
	ctx context.Context,
	queueJob model.QueueJob,
) error {
	logrus.WithFields(logrus.Fields{
		"queue_job": queueJob,
	}).Info("Executing questionnaire professional")

	return nil
}
