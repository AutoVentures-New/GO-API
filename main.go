package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/sirupsen/logrus"
	"github.com/trabalhe-conosco/api/config"
	"github.com/trabalhe-conosco/api/database"
	"github.com/trabalhe-conosco/api/log"
	"github.com/trabalhe-conosco/api/middleware"
	"github.com/trabalhe-conosco/api/router"
)

func main() {
	log.InitLogger()
	config.InitConfig()

	database.InitDatabase()

	database.RunMigrations()

	app := fiber.New(fiber.Config{
		Prefork:                  false,
		CaseSensitive:            false,
		StrictRouting:            false,
		ServerHeader:             "*",
		AppName:                  "Trabalhe Conosco API",
		Immutable:                true,
		DisableStartupMessage:    true,
		ErrorHandler:             middleware.ErrorHandler(),
		EnableSplittingOnParsers: true,
		EnablePrintRoutes:        false,
	})

	app.Use(cors.New())

	router.SetupRoutes(app)

	logrus.Info("API stated with success")

	logrus.Fatal(app.Listen(":" + config.Config.Port))
}
