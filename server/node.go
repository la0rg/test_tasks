package server

import (
	"net"
	"time"

	"google.golang.org/grpc"

	"github.com/la0rg/test_tasks/cache"
	"github.com/la0rg/test_tasks/rpc"
	"github.com/la0rg/test_tasks/util"
	"github.com/la0rg/test_tasks/vector_clock"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

const (
	REPLICATION = 3
	WRITETO     = 2
	READFROM    = 2
)

const WEIGHT uint8 = 10

type Node struct {
	name     string
	mbrship  *Membership
	httpPort int
	cache    *cache.Cache
}

func NewNode(address string, iport int) (*Node, error) {
	membership := NewMembership(address)
	addr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return nil, errors.Wrap(err, "Not able to resolve node address")
	}
	// TODO: node weight should be defined based on memory size
	err = membership.AddNode(address, iport, WEIGHT)
	if err != nil {
		return nil, err
	}
	cache := cache.NewCache()
	return &Node{
		name:     address,
		mbrship:  membership,
		httpPort: addr.Port,
		cache:    cache,
	}, nil
}

func (n *Node) CoordinatorPut(key string, value *cache.CacheValue, vc *vector_clock.VC) error {

	if vc == nil {
		// coordinator generates a vector clock for the new value
		vc = vector_clock.NewVc()
	}
	vc.Incr(n.name)
	// write localy
	n.cache.Set(key, value, vc)

	// The coordinator sends the new version (along with the new vector clock) to
	// the REPLICATION highest-ranked reachable nodes. If at least WRITETO-1 nodes respond then the write is considered successful.
	preferenceList := n.mbrship.FindPreferenceList(key, REPLICATION)
	done := make(chan struct{})
	// Replicate to REPLICATION nodes
	for _, endpoint := range preferenceList {

		// do not replicate on itself
		if endpoint.Address.String() == n.name {
			continue
		}

		log.Debugf("Send Set to replica: %s", endpoint.Address.String())
		// request Set on replica node
		go func(endpoint *Endpoint) {
			putRequest(endpoint)
			done <- struct{}{}
		}(endpoint)
	}

	timeout := util.WaitOnChanWithTimeout(done, WRITETO-1, time.Second)
	if timeout {
		return errors.New("Unsufficient number of replicas responed within a second")
	}
	return nil
}

func putRequest(endpoint *Endpoint) {
	conn, err := grpc.Dial(endpoint.IAddress(), []grpc.DialOption{grpc.WithInsecure()}...)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	client := rpc.NewNodeServiceClient(conn)
	_, err = client.Set(context.Background(), &rpc.ClockedValue{})
	if err != nil {
		log.Fatalf("Replica Set goes wrong...%v", err)
	}
	log.Debugf("Got result from replica: %s", endpoint.Address.String())
}

func (n *Node) Set(ctx context.Context, value *rpc.ClockedValue) (*rpc.SetResult, error) {
	log.Debug("Set method was called")
	return &rpc.SetResult{}, nil

	// Read local value
	// If local Version Vector descends incoming Version Vector ignore write (you’ve seen it!)
	// If Incoming Version Vector descends local Version Vector overwrite local value with new one
	// If Incoming Version Vector is concurrent with local Version Vector, merge values
}
