package server

import (
	"net"

	"github.com/la0rg/test_tasks/cache"
	"github.com/pkg/errors"
)

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
	err = membership.AddNode(address, iport)
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
