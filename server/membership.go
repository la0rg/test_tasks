package server

import (
	"fmt"
	"math/rand"
	"net"

	"github.com/la0rg/test_tasks/hash"
	"github.com/la0rg/test_tasks/rpc"
	"github.com/la0rg/test_tasks/vector_clock"
	"github.com/pkg/errors"
	"github.com/rs/xid"
	log "github.com/sirupsen/logrus"
)

type Membership struct {
	// Name contains address of the membership (current node)
	// Can be used to check if an endpoint is on the same machine
	Name string

	// Partitioning and replication ranges are based on the Ring
	ring hash.Ring

	// State that is gonna be transfferd with gossip
	Endpoints []*Endpoint

	// Several virtual nodes (virtual ids) may lead to the same Endpoint
	// VNodes should contains pointers to the same instances as in the endpoints
	VNodes map[string]*Endpoint
	Vc     *vector_clock.VC
}

func NewMembership(name string) *Membership {
	return &Membership{
		Name:      name,
		Endpoints: make([]*Endpoint, 0),
		VNodes:    make(map[string]*Endpoint, 0),
		Vc:        vector_clock.NewVc(),
	}
}

func (m *Membership) addEndpoint(e *Endpoint) {
	m.Endpoints = append(m.Endpoints, e)
	m.Vc.Incr(m.Name)
}

type Endpoint struct {
	// Public http address
	Address net.TCPAddr
	// Inter-node communication port
	IPort int

	// TODO: availability
}

func CompareEndpoints(e1, e2 Endpoint) bool {
	if &e1 == &e2 {
		return true
	}
	if e1.Address.Port != e2.Address.Port {
		return false
	}
	if !e1.Address.IP.Equal(e2.Address.IP) {
		return false
	}
	if e1.IPort != e2.IPort {
		return false
	}
	return true
}

func (e *Endpoint) IAddress() string {
	return fmt.Sprintf("%s:%d", e.Address.IP, e.IPort)
}

func (m *Membership) RndLiveEndpoint() *Endpoint {
	available := make([]int, 0)
	// filtering
	for i := range m.Endpoints {
		if m.Endpoints[i].Address.String() != m.Name {
			available = append(available, i)
		}
	}

	if len(available) > 0 {
		return m.Endpoints[available[rand.Intn(len(available))]]
	} else {
		return nil
	}
}

func (m *Membership) AddNode(name string, iport int, weight uint8) error {
	addr, err := net.ResolveTCPAddr("tcp", name)
	if err != nil {
		return errors.Wrap(err, "Not able to resolve node address")
	}
	endpoint := &Endpoint{
		Address: *addr,
		IPort:   iport,
	}
	m.addEndpoint(endpoint)
	for i := 0; i < int(weight); i++ {
		virtualId := xid.New().String()
		m.VNodes[virtualId] = endpoint
		m.ring.AddNode(virtualId)
	}
	return nil
}

func (m *Membership) ToRpc() *rpc.Membership {
	res := &rpc.Membership{
		Endpoints: make([]*rpc.Membership_Endpoint, len(m.Endpoints)),
		Vnodes:    make(map[string]*rpc.Membership_Endpoint, len(m.VNodes)),
		VectorClock: &rpc.VC{
			Store: m.Vc.GetStore(),
		},
	}

	for i, v := range m.Endpoints {
		res.Endpoints[i] = &rpc.Membership_Endpoint{
			Ip:    v.Address.IP,
			Port:  int32(v.Address.Port),
			Iport: int32(v.IPort),
		}
	}

	for k, v := range m.VNodes {
		res.Vnodes[k] = &rpc.Membership_Endpoint{
			Ip:    v.Address.IP,
			Port:  int32(v.Address.Port),
			Iport: int32(v.IPort),
		}
	}

	return res
}

func (m *Membership) MergeRpc(rpcMbr *rpc.Membership) {
	vc := &vector_clock.VC{Store: rpcMbr.VectorClock.Store}
	endpoints := make([]*Endpoint, 0)
	for _, endpoint := range rpcMbr.Endpoints {
		endpoints = append(endpoints, &Endpoint{
			Address: net.TCPAddr{IP: endpoint.GetIp(), Port: int(endpoint.GetPort())},
			IPort:   int(endpoint.GetIport()),
		})
	}
	vnodes := make(map[string]Endpoint, len(rpcMbr.Vnodes))
	for key, endpoint := range rpcMbr.Vnodes {
		vnodes[key] = Endpoint{
			Address: net.TCPAddr{IP: endpoint.GetIp(), Port: int(endpoint.GetPort())},
			IPort:   int(endpoint.GetIport()),
		}
	}
	switch vector_clock.Compare(m.Vc, vc) {
	case -1:
		// if VC of current node is staled from rpc VC then use rpc Membership as node membership
		m.Endpoints = endpoints
		m.VNodes = m.toMembershipEndpoints(vnodes)
		m.Vc = vc
		m.ring.Clear()
		for key := range m.VNodes {
			m.ring.AddNode(key)
		}
	case 0:
		if vector_clock.Equal(m.Vc, vc) {
			return
		}
		// if VCs are not comparible than merge nodes of both memberships
		for _, endpoint := range endpoints {
			m.Endpoints = append(m.Endpoints, endpoint)
		}
		for key, endpoint := range m.toMembershipEndpoints(vnodes) {
			m.VNodes[key] = endpoint
			m.ring.AddNode(key)
		}
		new_vc := vector_clock.Merge(m.Vc, vc)
		new_vc.Incr(m.Name)
		m.Vc = new_vc
	}
}

func Compare(mbr1, mbr2 Membership) bool {
	if &mbr1 == &mbr2 {
		return true
	}
	if len(mbr1.Endpoints) != len(mbr2.Endpoints) {
		return false
	}
	if len(mbr1.VNodes) != len(mbr2.VNodes) {
		return false
	}
	for i, v := range mbr1.Endpoints {
		if mbr2.Endpoints[i] != v {
			return false
		}
	}
	for k, v := range mbr1.VNodes {
		if mbr2.VNodes[k] != v {
			return false
		}
	}
	return true
}

func (m *Membership) toMembershipEndpoints(vnodes map[string]Endpoint) map[string]*Endpoint {
	updated := make(map[string]*Endpoint, len(vnodes))
l1:
	for k, e := range vnodes {
		for i := range m.Endpoints {
			if CompareEndpoints(*m.Endpoints[i], e) {
				updated[k] = m.Endpoints[i]
				continue l1
			}
		}
	}
	return updated
}

// The list of nodes that is responsible for storing a particular key.
// To account for node failures, preference list contains more
// than N nodes.
// Preference list for a key is constructed by skipping positions in the
// ring to ensure that the list contains only distinct physical nodes.
func (m *Membership) FindPreferenceList(key string, replication int) []*Endpoint {
	set := make(map[*Endpoint]bool, 0)
	preferenceList := make([]*Endpoint, 0)

	node := m.ring.FindNode(key)
	for len(preferenceList) < replication+1 {
		// translate virtual node name to physical endpoint
		endpoint := m.VNodes[node.GetName()]
		if _, ok := set[endpoint]; !ok {
			set[endpoint] = true
			preferenceList = append(preferenceList, endpoint)
		}
		node = m.ring.FindSuccessorNode(node)
	}
	return preferenceList
}

func (m *Membership) FindCoordinatorEndpoint(key string) *Endpoint {
	vnodeName := m.ring.FindNode(key).GetName()
	endpoint, ok := m.VNodes[vnodeName]
	if !ok {
		log.Fatalf("Internal state is corrupted: could not find endpoint by vnode name: %s", vnodeName)
	}
	return endpoint
}
