package server

import (
	"context"
	"net"
	"net/http"
	"strconv"
	"time"

	"encoding/json"

	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type Node struct {
	// TODO: cache etc.
	httpServer *http.Server
	mbrship    *Membership
	httpPort   int
}

func NewNode(address string) (*Node, error) {
	membership := NewMembership(address)
	// TODO: identify node addr on start (or programmatically)
	addr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return nil, errors.Wrap(err, "Not able to resolve node address")
	}
	err = membership.AddNode(address)
	if err != nil {
		return nil, err
	}
	return &Node{
		mbrship:  membership,
		httpPort: addr.Port,
	}, nil
}

func (n *Node) setupRouting(r *httprouter.Router) {
	// cache
	r.GET("/", n.Get)
	r.PUT("/", n.Update)
	r.POST("/", n.Set)
	r.DELETE("/", n.Remove)

	// configuration
	r.GET("/membership", n.Membership)
	r.GET("/membership/endpoint", n.Endpoint)
	r.POST("/membership/endpoint", n.AddEndpoint)
}

func (n *Node) StartHttpServer() {
	addr := ":" + strconv.Itoa(n.httpPort)
	router := httprouter.New()
	n.setupRouting(router)
	n.httpServer = &http.Server{
		Addr:         addr,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      router,
	}

	go func() {
		log.Infof("Start listening on %s", addr)
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

func (n *Node) Get(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	log.Debug("Processing get request")
}

func (n *Node) Set(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	log.Debug("Processing set request")
}

func (n *Node) Update(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	log.Debug("Processing update request")
}

func (n *Node) Remove(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	log.Debug("Processing remove request")
}

func (n *Node) Membership(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	v, err := json.Marshal(n.mbrship)
	if err != nil {
		log.Error(errors.Wrap(err, "Problem while marshaling membership struct"))
	}
	w.Write(v)
}

func (n *Node) Endpoint(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	v, err := json.Marshal(n.mbrship.Endpoints)
	if err != nil {
		log.Error(errors.Wrap(err, "Problem while marshaling membership struct"))
	}
	w.Write(v)
}

func (n *Node) AddEndpoint(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	decoder := json.NewDecoder(r.Body)
	var endpoints []Endpoint
	err := decoder.Decode(&endpoints)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()
	resp := NewRestResponse()
	for _, endpoint := range endpoints {
		err := n.mbrship.AddNode(endpoint.Address.String())
		if err != nil {
			resp.Error(err)
		}
	}
	w.Write(resp.Build())
}
