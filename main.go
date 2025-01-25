package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/hubjob/api/app/adapters/sendgrid"
	"github.com/hubjob/api/config"
	"github.com/hubjob/api/database"
	"github.com/hubjob/api/log"
	"github.com/hubjob/api/middleware"
	"github.com/hubjob/api/router"
	"github.com/sirupsen/logrus"
)

func main() {
	log.InitLogger()
	config.InitConfig()
	database.InitDatabase()
	database.RunMigrations()
	sendgrid.InitSendGrid()

	app := fiber.New(fiber.Config{
		Prefork:                  false,
		CaseSensitive:            false,
		StrictRouting:            false,
		ServerHeader:             "*",
		AppName:                  "HubJob API",
		Immutable:                true,
		DisableStartupMessage:    true,
		ErrorHandler:             middleware.ErrorHandler(),
		EnableSplittingOnParsers: true,
		EnablePrintRoutes:        false,
	})

	app.Use(cors.New())

	router.SetupRoutes(app)

	app.Use(router.RouteNotFound())

	logrus.Info("API stated with success")

	go func() {
		if err := app.Listen(":" + config.Config.Port); err != nil {
			logrus.WithError(err).Panic("Error on listen server")
		}
	}()

	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	_ = <-c

	logrus.Info("API stated graceful shutdown")

	_ = app.Shutdown()

	database.CloseDatabase()

	logrus.Info("API finish with sucess")
}
