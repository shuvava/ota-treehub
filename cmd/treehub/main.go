package main

import (
	"fmt"

	"github.com/shuvava/go-logging/logger"
	"github.com/sirupsen/logrus"

	"github.com/shuvava/treehub/pkg/version"

	"github.com/shuvava/treehub/internal/app"
)

func main() {
	log := logger.NewLogrusLogger(logrus.InfoLevel)
	log.Info(fmt.Sprintf("Starting %s/%s", version.AppName, version.Version))

	server := app.NewServer(log)
	server.Start()
}
