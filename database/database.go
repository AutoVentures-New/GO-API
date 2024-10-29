package database

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
	"github.com/trabalhe-conosco/api/config"
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
