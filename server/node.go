package server

import (
	"log"
	"net"

	"github.com/la0rg/test_tasks/cache"
	"github.com/la0rg/test_tasks/vector_clock"
	"github.com/pkg/errors"
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

// coordinator initiated write
func (n *Node) Put(key string, value *cache.CacheValue) {
	// coordinator generates a vector clock for the new value
	newVc := vector_clock.NewVc()
	newVc.Incr(n.name)
	// write localy
	n.cache.Set(key, value, newVc)
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
