package job

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/hubjob/api/config"
	"github.com/hubjob/api/pkg"
	"github.com/sirupsen/logrus"
	"mime/multipart"
)

func SaveJobVideo(
	ctx context.Context,
	companyID int64,
	video *multipart.FileHeader,
	contentType string,
) (string, error) {
	fileReader, err := video.Open()
	if err != nil {
		logrus.WithError(err).Error("Error to open video file")

		return "", err
	}

	defer fileReader.Close()

	videoPath := fmt.Sprintf(
		"job_videos/%d/%s",
		companyID,
		video.Filename,
	)

	result, err := pkg.S3Uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(config.Config.S3.BucketPublic),
		Key:         aws.String(videoPath),
		Body:        fileReader,
		ContentType: aws.String(contentType),
		ACL:         aws.String("public-read"),
	})
	if err != nil {
		logrus.WithError(err).Error("Error to upload video file")

		return "", err
	}

	return result.Location, nil
}
