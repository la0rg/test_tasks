package server

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type HttpServer struct {
	httpServer *http.Server
	*Node
}

func NewHttpServer(n *Node) *HttpServer {
	return &HttpServer{
		Node: n,
	}
}

func (s *HttpServer) setupRouting(r *httprouter.Router) {
	// cache
	r.GET("/", s.Get)
	r.PUT("/", s.Update)
	r.POST("/", s.Set)
	r.DELETE("/", s.Remove)

	// configuration
	r.GET("/membership", s.Membership)
	r.GET("/membership/endpoint", s.Endpoint)
	r.POST("/membership/endpoint", s.AddEndpoint)
}

func (s *HttpServer) Start() error {
	addr := ":" + strconv.Itoa(s.httpPort)
	router := httprouter.New()
	s.setupRouting(router)
	s.httpServer = &http.Server{
		Addr:         addr,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      router,
	}

	go func() {
		log.Infof("Start listening on %s", addr)
		if err := s.httpServer.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()
	return nil
}

func (s *HttpServer) Stop() error {
	log.Infof("Shutting down the http server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := s.httpServer.Shutdown(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (s *HttpServer) Get(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	log.Debug("Processing get request")
}

func (s *HttpServer) Set(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	log.Debug("Processing set request")
}

func (s *HttpServer) Update(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	log.Debug("Processing update request")
}

func (s *HttpServer) Remove(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	log.Debug("Processing remove request")
}

func (s *HttpServer) Membership(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	v, err := json.Marshal(s.mbrship)
	if err != nil {
		log.Error(errors.Wrap(err, "Problem while marshaling membership struct"))
	}
	w.Write(v)
}

func (s *HttpServer) Endpoint(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	v, err := json.Marshal(s.mbrship.Endpoints)
	if err != nil {
		log.Error(errors.Wrap(err, "Problem while marshaling membership struct"))
	}
	w.Write(v)
}

func (s *HttpServer) AddEndpoint(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	decoder := json.NewDecoder(r.Body)
	var endpoints []Endpoint
	err := decoder.Decode(&endpoints)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()
	resp := NewRestResponse()
	for _, endpoint := range endpoints {
		err := s.mbrship.AddNode(endpoint.Address.String(), endpoint.IPort)
		if err != nil {
			resp.Error(err)
		}
	}
	w.Write(resp.Build())
}
