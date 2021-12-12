package main

import (
	"github.com/willie68/cel-service/internal/logging"
	log "github.com/willie68/cel-service/internal/logging"
)

func init() {
	initLogging()
}

func main() {
	log.Logger.Info("starting")

	log.Logger.Info("finished")
}

func initLogging() {
	log.Logger.SetLevel(logging.Debug)
	log.Logger.Init()
}
