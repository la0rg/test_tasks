package main

import (
	"flag"
	"math/rand"
	"os"
	"os/signal"
	"time"

	"github.com/la0rg/test_tasks/server"
	log "github.com/sirupsen/logrus"
)

func setup() {
	log.SetFormatter(&log.TextFormatter{ForceColors: true})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel) // TODO: change to approptiate log level
	rand.Seed(time.Now().Unix())
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
	node, err := server.NewNode(*addr, *port)
	fatalErr(err)

	// http
	httpServer := server.NewHttpServer(node)
	err = httpServer.Start()
	fatalErr(err)

	// gossip
	gossipServer := server.NewGossipServer(node)
	err = gossipServer.Seed(*seed)
	fatalErr(err)
	err = gossipServer.Start(*port)
	fatalErr(err)

	// stop
	<-stop
	log.Info("Shutting down the node..")
	httpServer.Stop()
	gossipServer.Stop()
}

func fatalErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
