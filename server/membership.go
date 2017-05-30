package server

import (
	"fmt"
	"math/rand"
	"net"
	"strconv"

	"github.com/la0rg/test_tasks/hash"
	"github.com/la0rg/test_tasks/rpc"
	"github.com/la0rg/test_tasks/vector_clock"
	"github.com/pkg/errors"
)

type Membership struct {
	Name      string
	ring      hash.Ring  // partitioning and replication ranges are based on the Ring
	Endpoints []Endpoint // the state that is gonna be transfferd with gossip
	Vc        *vector_clock.VC
}

func NewMembership(name string) *Membership {
	return &Membership{
		Name:      name,
		Endpoints: make([]Endpoint, 0),
		Vc:        vector_clock.NewVc(),
	}
}

func (m *Membership) addEndpoint(e *Endpoint) {
	m.Endpoints = append(m.Endpoints, *e)
	m.Vc.Incr(m.Name)
}

type Endpoint struct {
	// Public http address
	Address net.TCPAddr
	// Inter-node communication port
	IPort int
}

func (e *Endpoint) IAddress() string {
	return fmt.Sprintf("%s:%d", e.Address.IP, e.IPort)
}

func (m *Membership) RndLiveEndpoint() *Endpoint {
	available := make([]*Endpoint, 0)
	// filtering
	for i := range m.Endpoints {
		if m.Endpoints[i].Address.String() != m.Name {
			available = append(available, &m.Endpoints[i])
		}
	}

	return available[rand.Intn(len(available))]
}

func (m *Membership) AddNode(name string, iport int) error {
	addr, err := net.ResolveTCPAddr("tcp", name)
	if err != nil {
		return errors.Wrap(err, "Not able to resolve node address")
	}
	m.addEndpoint(&Endpoint{Address: *addr, IPort: iport})
	m.ring.AddNode(name)
	return nil
}

func (m *Membership) ToRpc() *rpc.Membership {
	res := &rpc.Membership{
		Endpoints: make([]*rpc.Membership_Endpoint, 0),
		VectorClock: &rpc.VC{
			Store: m.Vc.GetStore(),
		},
	}

	for _, v := range m.Endpoints {
		endpoint := &rpc.Membership_Endpoint{Ip: v.Address.IP, Port: int32(v.Address.Port), Iport: int32(v.IPort)}
		res.Endpoints = append(res.Endpoints, endpoint)
	}

	return res
}

func (m *Membership) MergeRpc(rpcMbr *rpc.Membership) {
	vc := &vector_clock.VC{Store: rpcMbr.VectorClock.Store}
	endpoints := make([]Endpoint, 0)
	for _, endpoint := range rpcMbr.Endpoints {
		endpoints = append(endpoints, Endpoint{net.TCPAddr{IP: endpoint.GetIp(), Port: int(endpoint.GetPort())}, int(endpoint.GetIport())})
	}
	switch vector_clock.Compare(m.Vc, vc) {
	case -1:
		// if VC of current node is staled from rpc VC then use rpc Membership as node membership
		m.Endpoints = endpoints
		m.Vc = vc
		m.ring.Clear()
		for _, endpoint := range endpoints {
			m.ring.AddNode(string(endpoint.Address.IP) + ":" + strconv.Itoa(endpoint.Address.Port))
		}
	case 0:
		if vector_clock.Equal(m.Vc, vc) {
			return
		}
		// if VCs are not comparible than merge nodes of both memberships
		for _, endpoint := range endpoints {
			m.Endpoints = append(m.Endpoints, endpoint)
			m.ring.AddNode(string(endpoint.Address.IP) + ":" + strconv.Itoa(endpoint.Address.Port))
		}
		new_vc := vector_clock.Merge(m.Vc, vc)
		new_vc.Incr(m.Name)
		m.Vc = new_vc
	}
}
