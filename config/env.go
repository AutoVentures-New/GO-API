package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

var Config configEnv

type configEnv struct {
	Env                string `env:"ENV" envDefault:"prod"`
	FrontendURL        string `env:"FRONTEND_URL" envDefault:"http://localhost:3000"`
	Port               string `env:"PORT" envDefault:"5000"`
	NewRelicLicenseKey string `env:"NEW_RELIC_LICENSE_KEY" envDefault:"-"`
	JwtSecret          string `env:"JWT_SECRET"`
	Database           struct {
		Uri string `env:"DB_URI" json:"-"`
	}

	S3 struct {
		Region    string `env:"S3_REGION" json:"region"`
		Key       string `env:"S3_KEY" json:"-"`
		Secret    string `env:"S3_SECRET" json:"-"`
		Bucket    string `env:"S3_BUCKET" json:"bucket"`
		SaveFiles string `env:"S3_SAVE_FILES" json:"saveFiles"`
	}

	Redis struct {
		Address         string `env:"REDIS_ADDRESS"`
		Password        string `env:"REDIS_PASSWORD"`
		SessionDatabase int    `env:"REDIS_SESSION_DATABASE"`
	}
}

func InitConfig() {
	err := godotenv.Load()
	if err != nil {
		logrus.WithError(err).Error("Error loading .env file")
	}

	err = env.Parse(&Config)
	if err != nil {
		logrus.WithError(err).Fatal("error parsing config")
	}
}
