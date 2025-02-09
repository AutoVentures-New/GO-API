package pkg

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/hubjob/api/config"
	"github.com/sirupsen/logrus"
)

var SessionClient *redis.Client

func InitRedis() {
	SessionClient = redis.NewClient(&redis.Options{
		Addr:     config.Config.Redis.Address,
		Password: config.Config.Redis.Password,
		DB:       config.Config.Redis.SessionDatabase,
	})

	_, err := SessionClient.Ping(context.Background()).Result()
	if err != nil {
		logrus.WithError(err).Fatal("error to connect redis")

		return
	}

	logrus.Info("Successfully connect to redis")
}

func CloseRedis() {
	_ = SessionClient.Close()

	logrus.Info("Successfully close redis")
}
