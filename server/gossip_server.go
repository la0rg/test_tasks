package server

import (
	"fmt"
	"net"

	"github.com/la0rg/test_tasks/rpc"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type GossipServer struct {
	*Node
}

func NewGossipServer(n *Node) *GossipServer {
	return &GossipServer{
		Node: n,
	}
}

func (s GossipServer) ReqForMembership(ctx context.Context, in *rpc.Membership) (*rpc.Membership, error) {
	if s.mbrship == nil {
		return nil, errors.New("Gossip server started with nil membership")
	}

	// node requests membership for the first time and pass its membership to the seed node
	if endpoints := in.GetEndpoints(); len(endpoints) == 1 {
		s.mbrship.MergeRpc(in)
	}

	log.Debug("Answering for ReqForMembership")
	return s.mbrship.ToRpc(), nil
}

func (s GossipServer) Seed(address string) error {
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
		membership, err := client.ReqForMembership(context.Background(), s.mbrship.ToRpc())
		if err != nil {
			log.Error(errors.Wrap(err, "Problems on seeding round"))
		}
		log.Printf("membership = %+v\n", membership)
		s.mbrship.MergeRpc(membership)
	}()
	return nil
}

func (s GossipServer) Start(port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	grpcServer := grpc.NewServer()
	rpc.RegisterGossipServiceServer(grpcServer, s)
	log.Debugf("Start listening gossip on %d", port)
	go grpcServer.Serve(lis)
	return nil
}
