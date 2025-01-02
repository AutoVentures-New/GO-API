package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/sirupsen/logrus"
)

var Config configEnv

type configEnv struct {
	Env  string `env:"ENV" envDefault:"prod"`
	Port string `env:"PORT" envDefault:"5000"`

	JwtSecret string `env:"JWT_SECRET"`

	Database struct {
		Uri string `env:"DB_URI"`
	}

	SendGrid struct {
		ApiKey   string `env:"SENDGRID_API_KEY"`
		Sender   string `env:"SENDGRID_SENDER"`
		EmailDev string `env:"SENDGRID_EMAIL_DEV"`
	}
}

func InitConfig() {
	err := env.Parse(&Config)
	if err != nil {
		logrus.WithError(err).Fatal("error parsing config")
	}
}
