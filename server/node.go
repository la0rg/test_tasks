package server

import (
	"net"

	"github.com/la0rg/test_tasks/cache"
	"github.com/la0rg/test_tasks/vector_clock"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
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

func (n *Node) CoordinatorPut(key string, value *cache.CacheValue, vc *vector_clock.VC) {

	if vc == nil {
		// coordinator generates a vector clock for the new value
		vc = vector_clock.NewVc()
	}
	// increment VC only on coordinator node
	vc.Incr(n.name)
	// write localy
	n.cache.Set(key, value, vc)
	// write to WRITETO - 1 nodes
	preferenceList := n.mbrship.FindPreferenceList(key, REPLICATION)
	for _, endpoint := range preferenceList {
		log.Printf("endpoint = %+v\n", endpoint)
	}
	//for _, nodeName := range preferenceList {
	//endpoint := n.mbrship.findEndpointByNodeName(nodeName)
	//if endpoint == nil {
	//// TODO: return err
	//}
	//// TODO: add filtration based on node availability
	//go {
	//// RPC CALL TO ENDPOINT FOR WRITE
	//}()
	//}

	// wait for WRITETO - 1 nodes
}

func (n *Node) put(key string, value *cache.CacheValue, vc *vector_clock.VC) {
	// Read local value
	// If local Version Vector descends incoming Version Vector ignore write (you’ve seen it!)
	// If Incoming Version Vector descends local Version Vector overwrite local value with new one
	// If Incoming Version Vector is concurrent with local Version Vector, merge values
	clockedValue, ok := n.cache.Get(key)
	if ok {
		localVc := clockedValue.VC
		switch vector_clock.Compare(localVc, vc) {
		case 1: 
			return // ignore
		case 0:
			if vector_clock.Equal(localVc, vc) {
				break // overwrite value with equal vcs
			}
			// merging
			switch 
		}
	}
	n.cache.Set(key, value, vc)
}

//func (n *Node) PutValue(ctx context.Context, in *GCacheValue, opts ...grpc.CallOption) (*GError, error)
