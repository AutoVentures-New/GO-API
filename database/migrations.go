package database

import (
	"errors"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/sirupsen/logrus"
	"github.com/trabalhe-conosco/api/config"
)

func RunMigrations() {
	if config.Config.Env == "prod" {
		return
	}

	driver, err := mysql.WithInstance(Database, &mysql.Config{})
	if err != nil {
		logrus.WithError(err).Fatal("Could not connect to database")
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://./database/migrations",
		"trabalhe-conosco",
		driver,
	)
	if err != nil {
		logrus.WithError(err).Fatal("Could not connect to database")
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		logrus.WithError(err).Fatal("Could not run migration")
	}

	logrus.Info("Migrations run successfully")
}
