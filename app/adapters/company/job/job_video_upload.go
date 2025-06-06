package job

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

func SaveJobVideo(
	ctx context.Context,
	companyID int64,
	video *multipart.FileHeader,
	contentType string,
) (model.JobVideo, error) {
	fileReader, err := video.Open()
	if err != nil {
		logrus.WithError(err).Error("Error to open video file")

		return model.JobVideo{}, err
	}

	defer fileReader.Close()

	videoPath := fmt.Sprintf(
		"public/job_videos/%d/%s",
		companyID,
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

		return model.JobVideo{}, err
	}

	// 2. Gera URL CloudFront
	cloudfrontDomain := "d3437slqdk74db.cloudfront.net" // Substitua pelo seu
	publicURL := fmt.Sprintf("https://%s/%s", cloudfrontDomain, videoPath)

	jobVideo := model.JobVideo{
		CompanyID: companyID,
		VideoLink: publicURL,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	dbTransaction, err := database.Database.Begin()
	if err != nil {
		logrus.WithError(err).Error("Error to open db transaction")

		return model.JobVideo{}, err
	}

	_, err = dbTransaction.ExecContext(
		ctx,
		`INSERT INTO job_videos(company_id,video_link,created_at,updated_at) VALUES(?,?,?,?)`,
		jobVideo.CompanyID,
		jobVideo.VideoLink,
		jobVideo.CreatedAt,
		jobVideo.UpdatedAt,
	)
	if err != nil {
		logrus.WithError(err).Error("Error to insert job video")

		_ = dbTransaction.Rollback()

		return model.JobVideo{}, err
	}

	if err := dbTransaction.Commit(); err != nil {
		_ = dbTransaction.Rollback()

		logrus.WithError(err).Error("Error to commit transaction")

		return model.JobVideo{}, err
	}

	return model.JobVideo{}, err
}
