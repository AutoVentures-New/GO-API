package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

var Config configEnv

type configEnv struct {
	Env       string `env:"ENV" envDefault:"prod"`
	Port      string `env:"PORT" envDefault:"5000"`
	JwtSecret string `env:"JWT_SECRET"`
	Database  struct {
		Uri string `env:"DB_URI"`
	}
	SendGrid struct {
		ApiKey   string `env:"SENDGRID_API_KEY"`
		Sender   string `env:"SENDGRID_SENDER"`
		EmailDev string `env:"SENDGRID_EMAIL_DEV"`
	}
	S3 struct {
		Bucket          string `env:"S3_BUCKET_NAME"`
		BucketRegion    string `env:"S3_BUCKET_REGION"`
		AccessKey       string `env:"S3_ACCESS_KEY"`
		SecretAccessKey string `env:"S3_SECRET_ACCESS_KEY"`
	}
	Redis struct {
		Address         string `env:"REDIS_ADDRESS"`
		Password        string `env:"REDIS_PASSWORD"`
		SessionDatabase int    `env:"REDIS_SESSION_DATABASE"`
	}
	GCP struct {
		CredentialsFile string `env:"GCP_CREDENTIALS_FILE"`
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
