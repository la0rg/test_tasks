package server

import (
	"net"

	"github.com/la0rg/test_tasks/hash"
	"github.com/la0rg/test_tasks/vector_clock"
	"github.com/pkg/errors"
)

type Membership struct {
	ring      hash.Ring  // partitioning and replication ranges are based on the Ring
	endpoints []Endpoint // the state that is gonna be transfferd with gossip
	vc        *vector_clock.VC
}

func NewMembership() *Membership {
	return &Membership{
		endpoints: make([]Endpoint, 1),
		vc:        vector_clock.NewVc(),
	}
}

func (m *Membership) addEndpoint(e *Endpoint) {
	m.endpoints = append(m.endpoints, *e)
}

type Endpoint struct {
	address net.Addr
}

func (m *Membership) AddNode(name string) error {
	addr, err := net.ResolveTCPAddr("tcp", name)
	if err != nil {
		errors.Wrap(err, "Not able to resolve node address. The node will not be added to the cluster.")
	}
	m.addEndpoint(&Endpoint{address: addr})
	m.ring.AddNode(name)
	return nil
}
