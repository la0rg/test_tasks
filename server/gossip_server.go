package server

import (
	"net"
	"time"

	"github.com/la0rg/test_tasks/rpc"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type GossipServer struct {
	*Node
	quiteChan chan struct{}
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
		log.Debug("Merging incoming membership from the connected new node")
	}

	log.Debug("Answering for ReqForMembership")
	return s.mbrship.ToRpc(), nil
}

func (s GossipServer) Seed(address string) error {
	err := s.gossipRequest(address, s.mbrship)
	if err != nil {
		return err
	}
	return nil
}

func (s GossipServer) StopGossipRoutine() {
	if s.quiteChan != nil {
		close(s.quiteChan)
	}
}

func (s GossipServer) StartGossipRoutine() {
	ticker := time.NewTicker(2 * time.Second)
	s.quiteChan = make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				// Gossip with a random live node
				if rndEndpoint := s.mbrship.RndLiveEndpoint(); rndEndpoint != nil {
					addr := rndEndpoint.IAddress()
					log.Debugf("Start gossip communication %s", addr)
					s.gossipRequest(addr, NewMembership(""))
				}
				// Gossip with a random dead node
				// Maybe gossip with a seed
			case <-s.quiteChan:
				ticker.Stop()
				return
			}
		}
	}()
}

func (s GossipServer) gossipRequest(address string, membership *Membership) error {
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
		membership, err := client.ReqForMembership(context.Background(), membership.ToRpc())
		if err != nil {
			log.Error(errors.Wrap(err, "Problems on seeding round"))
		}
		s.mbrship.MergeRpc(membership)
	}()
	return nil
}
