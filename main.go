package main

import (
	"flag"
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

	// flags
	addr := flag.String("addr", "", "Public addres (ip:port) of the current node.")
	flag.Parse()
	if *addr == "" {
		flag.Usage()
		log.Fatal("Public address cannot be empty.")
	}

	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)

	// init
	node, err := server.NewNode(*addr)
	if err != nil {
		log.Fatal(err)
	}
	node.StartHttpServer()

	// stop
	<-stop
	log.Info("Shutting down the node..")
	node.StopHttpServer()
}
