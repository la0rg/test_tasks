package server

import (
	"github.com/la0rg/test_tasks/rpc"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

type GossipServer struct {
	membership *Membership
}

func (g GossipServer) ReqForMembership(ctx context.Context, in *rpc.Membership) (*rpc.Membership, error) {
	if g.membership == nil {
		return nil, errors.New("Gossip server started with nil membership")
	}

	// node requests membership for the first time and pass its address to add to membership of the seed node
	if endpoints := in.GetEndpoints(); len(endpoints) == 1 {
		g.membership.MergeRpc(in)
	}

	log.Debug("Answering for ReqForMembership")
	return g.membership.ToRpc(), nil
}
