package config

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func LogInit(app AppConfig) {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetReportCaller(true)
	logLevel, _ := log.ParseLevel(app.LogLevel)
	log.SetLevel(logLevel)
	log.SetOutput(os.Stdout)
}
