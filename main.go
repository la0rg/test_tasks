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
	addr := flag.String("addr", "", "Public addres(ip:port) of the current node.")
	seed := flag.String("seed", "", "Address of the seed server.")
	port := flag.Int("internal_port", 7770, "Port for inter-node communication")
	flag.Parse()

	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)

	// init
	node, err := server.NewNode(*addr)
	fatalErr(err)
	err = node.StartHttpServer()
	fatalErr(err)
	err = node.Seed(*seed)
	fatalErr(err)
	err = node.StartGossipServer(*port) // TODO: hardcoded port

	// stop
	<-stop
	log.Info("Shutting down the node..")
	node.StopHttpServer()
}

func fatalErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
