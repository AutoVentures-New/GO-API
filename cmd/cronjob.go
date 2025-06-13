package cmd

import (
	"github.com/AutoVentures-New/GO-API/config"
	"github.com/AutoVentures-New/GO-API/database"
	"github.com/AutoVentures-New/GO-API/log"
	"github.com/AutoVentures-New/GO-API/pkg"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var cronjobCmd = &cobra.Command{
	Use:   "cronjob",
	Short: "Autoventures Cronjob",
	Run: func(cmd *cobra.Command, args []string) {
		log.InitLogger()
		config.InitConfig()
		database.InitDatabase()
		pkg.InitS3Client()

		database.CloseDatabase()

		logrus.Info("Cronjob finish with success")
	},
}

func init() {
	rootCmd.AddCommand(cronjobCmd)
}
