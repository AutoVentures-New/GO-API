package profile

import (
	"context"
	"fmt"
	"mime/multipart"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/hubjob/api/config"
	"github.com/hubjob/api/database"
	"github.com/hubjob/api/pkg"
	"github.com/sirupsen/logrus"
)

func UpdateCandidatePhoto(
	ctx context.Context,
	candidateID int64,
	photo *multipart.FileHeader,
	contentType string,
) error {
	fileReader, err := photo.Open()
	if err != nil {
		logrus.WithError(err).Error("Error to open photo file")

		return err
	}

	defer fileReader.Close()

	photoPath := fmt.Sprintf(
		"candidates/%d/photo/%s",
		candidateID,
		photo.Filename,
	)

	_, err = pkg.S3Uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(config.Config.S3.Bucket),
		Key:         aws.String(photoPath),
		Body:        fileReader,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		logrus.WithError(err).Error("Error to upload photo file")

		return err
	}

	dbTransaction, err := database.Database.Begin()
	if err != nil {
		logrus.WithError(err).Error("Error to open db transaction")

		return err
	}

	now := time.Now().UTC()

	_, err = dbTransaction.ExecContext(
		ctx,
		`INSERT INTO candidate_photos(candidate_id,bucket_name,photo_path,created_at,updated_at) 
				VALUES(?,?,?,?,?) 
				ON DUPLICATE KEY UPDATE
					bucket_name = VALUES(bucket_name),
					photo_path = VALUES(photo_path),
					updated_at = VALUES(updated_at)`,
		candidateID,
		config.Config.S3.Bucket,
		photoPath,
		now,
		now,
	)
	if err != nil {
		logrus.WithError(err).Error("Error to insert candidate photo")

		_ = dbTransaction.Rollback()

		return err
	}

	if err := dbTransaction.Commit(); err != nil {
		_ = dbTransaction.Rollback()

		logrus.WithError(err).Error("Error to commit transaction")

		return err
	}

	return nil
}
