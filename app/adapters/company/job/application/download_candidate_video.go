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

var ErrVideoNotFound = errors.New("video not found")

func DownloadCandidateVideo(
	ctx context.Context,
	applicationID int64,
) (*s3.GetObjectOutput, error) {
	var bucketName, videoPath string

	err := database.Database.QueryRowContext(
		ctx,
		`SELECT bucket_name,video_path 
				FROM job_application_candidate_videos WHERE application_id = ?`,
		applicationID,
	).Scan(
		&bucketName,
		&videoPath,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrVideoNotFound
	}

	if err != nil {
		logrus.WithError(err).Error("Error to get candidate video")

		return nil, err
	}

	input := &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(videoPath),
	}

	result, err := pkg.S3Client.GetObject(input)
	if err != nil {
		logrus.WithError(err).
			WithField("bucket", bucketName).
			WithField("video", videoPath).
			Error("Error to download video")

		return nil, ErrVideoNotFound
	}

	return result, nil
}
