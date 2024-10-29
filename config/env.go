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
}

func InitConfig() {
	err := env.Parse(&Config)
	if err != nil {
		logrus.WithError(err).Fatal("error parsing config")
	}
}
