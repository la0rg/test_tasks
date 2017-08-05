package server

import (
	"fmt"
	"net"

	"github.com/la0rg/test_tasks/rpc"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func StartGrpcServer(port int, gossipService rpc.GossipServiceServer, nodeService rpc.NodeServiceServer) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	grpcServer := grpc.NewServer()
	rpc.RegisterGossipServiceServer(grpcServer, gossipService)
	rpc.RegisterNodeServiceServer(grpcServer, nodeService)
	log.Debugf("Start listening internal communication on %d", port)
	go grpcServer.Serve(lis)
	return nil
}
