package server

import (
	"net"

	"github.com/pkg/errors"
)

type Node struct {
	// TODO: cache etc.
	name     string
	mbrship  *Membership
	httpPort int
}

func NewNode(address string, iport int) (*Node, error) {
	membership := NewMembership(address)
	// TODO: identify node addr on start (or programmatically)
	addr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return nil, errors.Wrap(err, "Not able to resolve node address")
	}
	err = membership.AddNode(address, iport)
	if err != nil {
		return nil, err
	}
	return &Node{
		name:     address,
		mbrship:  membership,
		httpPort: addr.Port,
	}, nil
}
