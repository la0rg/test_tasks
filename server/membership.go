package server

import (
	"net"

	"github.com/la0rg/test_tasks/hash"
	"github.com/la0rg/test_tasks/vector_clock"
	"github.com/pkg/errors"
)

type Membership struct {
	name      string
	ring      hash.Ring  // partitioning and replication ranges are based on the Ring
	Endpoints []Endpoint // the state that is gonna be transfferd with gossip
	Vc        *vector_clock.VC
}

func NewMembership(name string) *Membership {
	return &Membership{
		name:      name,
		Endpoints: make([]Endpoint, 0),
		Vc:        vector_clock.NewVc(),
	}
}

func (m *Membership) addEndpoint(e *Endpoint) {
	m.Endpoints = append(m.Endpoints, *e)
	m.Vc.Incr(m.name)
}

type Endpoint struct {
	Address net.TCPAddr
}

func (m *Membership) AddNode(name string) error {
	addr, err := net.ResolveTCPAddr("tcp", name)
	if err != nil {
		return errors.Wrap(err, "Not able to resolve node address")
	}
	m.addEndpoint(&Endpoint{Address: *addr})
	m.ring.AddNode(name)
	return nil
}
