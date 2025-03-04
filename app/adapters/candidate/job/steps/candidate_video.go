package steps

import (
	"context"
	"fmt"
	"mime/multipart"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/hubjob/api/config"
	"github.com/hubjob/api/database"
	"github.com/hubjob/api/model"
	"github.com/hubjob/api/pkg"
	"github.com/sirupsen/logrus"
)

func SaveCandidateVideo(
	ctx context.Context,
	application model.Application,
	video *multipart.FileHeader,
	contentType string,
) (model.Application, error) {
	fileReader, err := video.Open()
	if err != nil {
		logrus.WithError(err).Error("Error to open video file")

		return application, err
	}

	defer fileReader.Close()

	videoPath := fmt.Sprintf(
		"jobs/candidate_videos/%d/%d/%d/%s",
		application.CompanyID,
		application.JobID,
		application.ID,
		video.Filename,
	)

	_, err = pkg.S3Uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(config.Config.S3.Bucket),
		Key:         aws.String(videoPath),
		Body:        fileReader,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		logrus.WithError(err).Error("Error to upload video file")

		return application, err
	}

	dbTransaction, err := database.Database.Begin()
	if err != nil {
		logrus.WithError(err).Error("Error to open db transaction")

		return application, err
	}

	now := time.Now().UTC()
	candidateVideo := model.JobApplicationCandidateVideo{
		ApplicationID: application.ID,
		BucketName:    config.Config.S3.Bucket,
		VideoPath:     videoPath,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	_, err = dbTransaction.ExecContext(
		ctx,
		`INSERT INTO job_application_candidate_videos(application_id,bucket_name,video_path,score,created_at,updated_at) VALUES(?,?,?,?,?,?)`,
		candidateVideo.ApplicationID,
		candidateVideo.BucketName,
		candidateVideo.VideoPath,
		candidateVideo.Score,
		candidateVideo.CreatedAt,
		candidateVideo.UpdatedAt,
	)
	if err != nil {
		logrus.WithError(err).Error("Error to insert candidate video")

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
