package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"google.golang.org/grpc"

	"encoding/json"

	"github.com/julienschmidt/httprouter"
	"github.com/la0rg/test_tasks/rpc"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type Node struct {
	// TODO: cache etc.
	name         string
	httpServer   *http.Server
	gossipServer rpc.GossipServiceServer
	mbrship      *Membership
	httpPort     int
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
		name:     address,
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

func (n *Node) StartHttpServer() error {
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
	return nil
}

func (n *Node) StopHttpServer() error {
	log.Infof("Shutting down the http server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := n.httpServer.Shutdown(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (n *Node) StartGossipServer(port int) error {
	n.gossipServer = GossipServer{
		membership: n.mbrship,
	}
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	grpcServer := grpc.NewServer()
	rpc.RegisterGossipServiceServer(grpcServer, n.gossipServer)
	log.Debugf("Start listening gossip on %d", port)
	go grpcServer.Serve(lis)
	return nil
}

func (n *Node) Seed(address string) error {
	if address == "" {
		return nil
	}
	// validate address
	_, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return errors.Wrap(err, "Not able to resolve seed address")
	}
	// request membership status from the seed node
	go func() {
		conn, err := grpc.Dial(address, []grpc.DialOption{grpc.WithInsecure()}...)
		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()
		client := rpc.NewGossipServiceClient(conn)
		log.Debug("Start request for membership")
		membership, err := client.ReqForMembership(context.Background(), n.mbrship.ToRpc())
		if err != nil {
			log.Error(errors.Wrap(err, "Problems on seeding round"))
		}
		log.Printf("membership = %+v\n", membership)
		n.mbrship.MergeRpc(membership)
	}()
	return nil
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
