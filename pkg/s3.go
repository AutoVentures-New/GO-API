package pkg

import (
	"github.com/AutoVentures-New/GO-API/config"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/sirupsen/logrus"
)

var S3Uploader *s3manager.Uploader

var S3Client *s3.S3

func InitS3Client() {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(config.Config.S3.Region),
		Credentials: credentials.NewStaticCredentials(config.Config.S3.Key, config.Config.S3.Secret, ""),
	})
	if err != nil {
		logrus.WithError(err).Panic("Error to create s3 client")
	}

	S3Uploader = s3manager.NewUploader(sess)

	S3Client = s3.New(sess)
}
