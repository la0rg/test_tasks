package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/la0rg/test_tasks/util"
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
	// should not add endpoints manually (only connection via seed node)
	//r.POST("/membership/endpoint", s.AddEndpoint)
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
	key := r.URL.Query().Get("key")
	if len(key) == 0 {
		http.Error(w, "Parameter \"key\" should be specified", http.StatusBadRequest)
		return
	}
	value, ok := s.cache.Get(key)
	if !ok {
		http.Error(w, "No value for the specified key", http.StatusNotFound)
		return
	}
	ivalue, err := util.CacheValueToJson(value.CacheValue)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	v, err := json.Marshal(ivalue)
	if err != nil {
		err = errors.Wrap(err, "Problem while marshaling CacheValue struct")
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(v)
}

type test_struct struct {
	Value string
}

func (s *HttpServer) Set(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	log.Debug("Processing set request")
	jsonBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	key, value, err := util.ParseJson(jsonBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Forward request to the coordinator node for the specified key
	coordinatorEndpoint := s.mbrship.FindCoordinatorEndpoint(key)

	// This machine is a coordinator
	if s.mbrship.Name == coordinatorEndpoint.Address.String() {
		log.Debugf("Set for the key: %s, following value: %v", key, value)
		err := s.CoordinatorPut(key, value, nil) // TODO: retrieve VC from the request
		if err != nil {
			status := http.StatusInternalServerError
			if err == ErrTimeout || err == ErrCancel {
				status = http.StatusBadRequest
			}
			http.Error(w, err.Error(), status)
		}
		return
	}

	// Forward request to coordinator via reverse-proxy
	url, err := url.Parse(fmt.Sprintf("http://%s/", coordinatorEndpoint.Address.String()))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	r.Body = ioutil.NopCloser(bytes.NewBuffer(jsonBody)) // reset r.Body to make a proxy forwarding
	proxy := httputil.NewSingleHostReverseProxy(url)
	log.Debugf("Fowarding to coordinator with url: %v", url)
	proxy.ServeHTTP(w, r)
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
