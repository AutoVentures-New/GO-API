package profile

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

var ErrPhotoNotFound = errors.New("photo not found")

func DownloadCandidatePhoto(
	ctx context.Context,
	candidateID int64,
) (*s3.GetObjectOutput, error) {
	var bucketName, photoPath string

	err := database.Database.QueryRowContext(
		ctx,
		`SELECT bucket_name,photo_path 
				FROM candidate_photos WHERE candidate_id = ?`,
		candidateID,
	).Scan(
		&bucketName,
		&photoPath,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrPhotoNotFound
	}

	if err != nil {
		logrus.WithError(err).Error("Error to get candidate photo")

		return nil, err
	}

	input := &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(photoPath),
	}

	result, err := pkg.S3Client.GetObject(input)
	if err != nil {
		logrus.WithError(err).
			WithField("bucket", bucketName).
			WithField("photo", photoPath).
			Error("Error to download photo")

		return nil, ErrPhotoNotFound
	}

	return result, nil
}
