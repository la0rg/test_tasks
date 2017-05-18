package server

import (
	"context"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	CONN_PORT = "8090"
	CONN_HOST = "0.0.0.0"
)

type Node struct {
	// TODO: cache etc.
	httpServer *http.Server
}

func (n *Node) StartHttpServer() {
	addr := CONN_HOST + ":" + CONN_PORT
	n.httpServer = &http.Server{
		Addr:         addr,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	http.HandleFunc("/", methodRouter)

	go func() {
		log.Infof("Start listening on: %s", addr)
		if err := n.httpServer.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

}

func (n *Node) StopHttpServer() {
	log.Infof("Shutting down the http server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := n.httpServer.Shutdown(ctx)
	if err != nil {
		log.Fatal(err)
	}
}

func methodRouter(w http.ResponseWriter, r *http.Request) {
	log.Debugf("Handle HTTP request with method: %s", r.Method)
	switch r.Method {
	case http.MethodGet:
		get(w, r)
	case http.MethodPost:
		set(w, r)
	case http.MethodPut:
		update(w, r)
	case http.MethodDelete:
		remove(w, r)
	default:
		http.Error(w, "Method is not supported", http.StatusMethodNotAllowed)
	}
}

func get(w http.ResponseWriter, r *http.Request) {
	log.Debug("Processing get request")
}

func set(w http.ResponseWriter, r *http.Request) {
	log.Debug("Processing set request")
}

func update(w http.ResponseWriter, r *http.Request) {
	log.Debug("Processing update request")
}

func remove(w http.ResponseWriter, r *http.Request) {
	log.Debug("Processing remove request")
}
