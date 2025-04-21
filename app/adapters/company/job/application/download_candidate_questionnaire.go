package job

import (
	"context"
	"database/sql"
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/hubjob/api/database"
	"github.com/hubjob/api/pkg"
	"github.com/sirupsen/logrus"
)

var ErrFileNotFound = errors.New("file not found")

func DownloadCandidateQuestionnaire(
	ctx context.Context,
	companyID int64,
	applicationID int64,
	questionnaireType string,
) (*s3.GetObjectOutput, error) {
	var candidateID int64

	err := database.Database.QueryRowContext(
		ctx,
		`SELECT candidate_id 
				FROM job_applications WHERE company_id = ? AND id = ?`,
		companyID,
		applicationID,
	).Scan(&candidateID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrApplicationNotFound
	}

	if err != nil {
		logrus.WithError(err).Error("Error to get job application")

		return nil, err
	}

	var bucketName, resultFilePath string

	err = database.Database.QueryRowContext(
		ctx,
		`SELECT bucket_name,result_file_path 
				FROM candidate_questionnaires WHERE candidate_id = ? and type = ?
				ORDER BY created_at DESC
				LIMIT 1`,
		candidateID,
		questionnaireType,
	).Scan(
		&bucketName,
		&resultFilePath,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrPhotoNotFound
	}

	if err != nil {
		logrus.WithError(err).Error("Error to get candidate questionnaire")

		return nil, err
	}

	input := &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(resultFilePath),
	}

	result, err := pkg.S3Client.GetObject(input)
	if err != nil {
		logrus.WithError(err).
			WithField("bucket", bucketName).
			WithField("result_file_path", resultFilePath).
			Error("Error to download result file")

		return nil, ErrFileNotFound
	}

	return result, nil
}
