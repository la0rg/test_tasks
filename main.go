package main

import (
	"os"
	"os/signal"

	"github.com/la0rg/test_tasks/server"
	log "github.com/sirupsen/logrus"
)

func setup() {
	log.SetFormatter(&log.TextFormatter{ForceColors: true})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel) // TODO: change to approptiate log level
}

func main() {
	setup()
	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)

	// init
	node := server.Node{}
	node.StartHttpServer()

	// stop
	<-stop
	log.Info("Shutting down the node..")
	node.StopHttpServer()
}
