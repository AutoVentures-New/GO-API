package cmd

import (
	"github.com/hubjob/api/config"
	"github.com/hubjob/api/database"
	"github.com/hubjob/api/handler/queue_job"
	"github.com/hubjob/api/log"
	"github.com/hubjob/api/pkg"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var cronjobCmd = &cobra.Command{
	Use:   "cronjob",
	Short: "Hubjob Cronjob",
	Run: func(cmd *cobra.Command, args []string) {
		log.InitLogger()
		config.InitConfig()
		database.InitDatabase()
		pkg.InitS3Client()

		queue_job.Executor(cmd.Context())

		database.CloseDatabase()

		logrus.Info("Cronjob finish with success")
	},
}

func init() {
	rootCmd.AddCommand(cronjobCmd)
}
