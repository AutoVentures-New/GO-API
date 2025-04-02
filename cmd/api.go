package cmd

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
	"github.com/hubjob/api/pkg"
	"github.com/hubjob/api/router"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Hubjob API",
	Run: func(cmd *cobra.Command, args []string) {
		log.InitLogger()
		config.InitConfig()
		database.InitDatabase()
		database.RunMigrations()
		sendgrid.InitSendGrid()
		pkg.InitS3Client()
		pkg.InitRedis()

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
		pkg.CloseRedis()

		logrus.Info("API finish with success")
	},
}

func init() {
	rootCmd.AddCommand(apiCmd)
}
