package cmd

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/AutoVentures-New/GO-API/config"
	"github.com/AutoVentures-New/GO-API/database"
	"github.com/AutoVentures-New/GO-API/log"
	"github.com/AutoVentures-New/GO-API/middleware"
	"github.com/AutoVentures-New/GO-API/pkg"
	"github.com/AutoVentures-New/GO-API/router"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/newrelic/go-agent/v3/integrations/nrfiber"
	"github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "AutoVentures API",
	Run: func(cmd *cobra.Command, args []string) {
		log.InitLogger()
		config.InitConfig()
		database.InitDatabase()
		pkg.InitS3Client()
		pkg.InitNewRelic()

		app := fiber.New(fiber.Config{
			Prefork:                  false,
			CaseSensitive:            false,
			StrictRouting:            false,
			ServerHeader:             "*",
			AppName:                  "AutoVentures API",
			Immutable:                true,
			DisableStartupMessage:    true,
			ErrorHandler:             middleware.ErrorHandler(),
			EnableSplittingOnParsers: true,
			EnablePrintRoutes:        false,
			BodyLimit:                20 * 1024 * 1024,
		})

		app.Use(nrfiber.Middleware(pkg.NewRelicApp))

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

		logrus.Info("API finish with success")
	},
}

func init() {
	rootCmd.AddCommand(apiCmd)
}
