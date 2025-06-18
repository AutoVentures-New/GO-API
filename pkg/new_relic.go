package pkg

import (
	"github.com/AutoVentures-New/GO-API/config"
	"github.com/newrelic/go-agent/v3/newrelic"
	"log"
	"time"
)

var NewRelicApp *newrelic.Application

func InitNewRelic() {
	app, err := newrelic.NewApplication(
		newrelic.ConfigAppName("Autoventures Go API"),
		newrelic.ConfigLicense(config.Config.NewRelicLicenseKey),
		newrelic.ConfigDistributedTracerEnabled(true),
	)
	if err != nil {
		log.Fatalln("Error to connect to New Relic:", err)
	}

	err = app.WaitForConnection(5 * time.Second)
	if err != nil {
		log.Fatalln("Error to connect to New Relic:", err)
		return
	}
	NewRelicApp = app
}
