package database

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/hubjob/api/config"
	"github.com/sirupsen/logrus"
)

var Database *sql.DB

func InitDatabase() {
	var err error

	Database, err = sql.Open("mysql", config.Config.Database.Uri)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to connect to database")
	}

	Database.SetMaxOpenConns(10)
	Database.SetMaxIdleConns(5)

	if err := Database.Ping(); err != nil {
		logrus.WithError(err).Fatal("Failed to ping database")
	}

	logrus.Info("Successfully connected to database")
}

func CloseDatabase() {
	if Database == nil {
		return
	}

	if err := Database.Close(); err != nil {
		logrus.WithError(err).Error("Failed to close database")

		return
	}

	logrus.Info("Successfully closed database")
}
