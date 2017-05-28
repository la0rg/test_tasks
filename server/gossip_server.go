package server

import (
	"github.com/la0rg/test_tasks/rpc"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

type GossipServer struct{}

func (g GossipServer) ReqForMembership(ctx context.Context, in *rpc.GossipRequest) (*rpc.Membership, error) {
	log.Debug("Answering for ReqForMembership")
	return &rpc.Membership{VectorClock: &rpc.VC{Store: map[string]uint64{"test": 1}}}, nil
}
