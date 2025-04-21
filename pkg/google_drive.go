package pkg

import (
	"context"

	"github.com/hubjob/api/config"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

var DriveService *drive.Service

var SheetsService *sheets.Service

func InitSheetsService(ctx context.Context) {
	var err error

	SheetsService, err = sheets.NewService(
		ctx,
		option.WithCredentialsFile(config.Config.GCP.CredentialsFile),
	)
	if err != nil {
		logrus.WithError(err).Panic("Error to create sheets service")
	}
}

func InitDriveService(ctx context.Context) {
	var err error

	DriveService, err = drive.NewService(
		ctx,
		option.WithCredentialsFile(config.Config.GCP.CredentialsFile),
	)
	if err != nil {
		logrus.WithError(err).Panic("Error to create drive service")
	}
}
