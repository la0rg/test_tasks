package vector_clock

import (
	"sync"
)

type VC struct {
	mx    sync.Mutex
	store map[string]uint64 // TODO: add timestamps
}

func NewVc() *VC {
	return &VC{
		store: make(map[string]uint64),
	}
}

func (vc *VC) Incr(node string) {
	vc.mx.Lock()
	vc.store[node] = vc.store[node] + 1
	vc.mx.Unlock()
}

func (vc *VC) Get(node string) uint64 {
	vc.mx.Lock()
	defer vc.mx.Unlock()
	return vc.store[node]
}

func (vc *VC) happensBefore(other *VC) bool {
	vc.mx.Lock()
	defer vc.mx.Unlock()
	other.mx.Lock()
	defer other.mx.Unlock()

	// all the values are less or equal and at least one node on the the other clock is bigger
	before := false
	for node := range allKeys(vc.store, other.store) {
		if vc.store[node] > other.store[node] {
			return false
		} else if vc.store[node] < other.store[node] {
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

	for node := range allKeys(vc1.store, vc2.store) {
		result.store[node] = max(vc1.store[node], vc2.store[node])
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
