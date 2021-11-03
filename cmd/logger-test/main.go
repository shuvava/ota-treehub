package main

import (
	"github.com/shuvava/treehub/internal/logger"
	"github.com/sirupsen/logrus"
)

func main() {
	log := logger.NewLogrusLogger(logrus.DebugLevel)
	log.Trace("Something very low level.")
	log.Debug("Useful debugging information.")
	log.Info("Something noteworthy happened!")
	log.Warn("You should probably take a look at this.")
	log.Error("Something failed but I'm not quitting.")
	// Calls os.Exit(1) after logging
	log.Fatal("Bye.")
}
