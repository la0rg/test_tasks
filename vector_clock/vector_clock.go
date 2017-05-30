package vector_clock

import (
	"sync"
)

type VC struct {
	mx    sync.Mutex
	Store map[string]uint64 // TODO: add timestamps
}

func NewVc() *VC {
	return &VC{
		Store: make(map[string]uint64),
	}
}

func (vc *VC) Incr(node string) {
	vc.mx.Lock()
	vc.Store[node] = vc.Store[node] + 1
	vc.mx.Unlock()
}

func (vc *VC) Get(node string) uint64 {
	vc.mx.Lock()
	defer vc.mx.Unlock()
	return vc.Store[node]
}

func (vc *VC) GetStore() map[string]uint64 {
	res := make(map[string]uint64)
	vc.mx.Lock()
	defer vc.mx.Unlock()
	for k, v := range vc.Store {
		res[k] = v
	}
	return res
}

func (vc *VC) happensBefore(other *VC) bool {
	vc.mx.Lock()
	defer vc.mx.Unlock()
	other.mx.Lock()
	defer other.mx.Unlock()

	// all the values are less or equal and at least one node on the the other clock is bigger
	before := false
	for node := range allKeys(vc.Store, other.Store) {
		if vc.Store[node] > other.Store[node] {
			return false
		} else if vc.Store[node] < other.Store[node] {
			before = true
		}
	}
	return before
}

func Merge(vc1, vc2 *VC) *VC {
	result := NewVc()
	vc1.mx.Lock()
	defer vc1.mx.Unlock()
	vc2.mx.Lock()
	defer vc2.mx.Unlock()

	for node := range allKeys(vc1.Store, vc2.Store) {
		result.Store[node] = max(vc1.Store[node], vc2.Store[node])
	}
	return result
}

func Compare(vc1, vc2 *VC) int {
	if vc1.happensBefore(vc2) {
		return -1
	} else if vc2.happensBefore(vc1) {
		return 1
	} else {
		return 0
	}
}

func Equal(vc1, vc2 *VC) bool {
	vc1.mx.Lock()
	defer vc1.mx.Unlock()
	vc2.mx.Lock()
	defer vc2.mx.Unlock()

	if len(vc1.Store) != len(vc2.Store) {
		return false
	}

	for k, v := range vc1.Store {
		if v2, ok := vc2.Store[k]; !ok || v != v2 {
			return false
		}
	}
	return true
}

func max(a, b uint64) uint64 {
	if a > b {
		return a
	}
	return b
}

func allKeys(a, b map[string]uint64) map[string]bool {
	result := make(map[string]bool, len(a)+len(b))
	for node := range a {
		result[node] = true
	}
	for node := range b {
		result[node] = true
	}
	return result
}
