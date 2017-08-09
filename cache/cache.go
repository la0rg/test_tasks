package cache

import (
	"sync"

	"github.com/la0rg/test_tasks/vector_clock"
	log "github.com/sirupsen/logrus"
)

type Cache struct {
	mx    sync.Mutex
	store map[string]ClockedValue
}

func NewCache() *Cache {
	return &Cache{store: make(map[string]ClockedValue)}
}

type CType uint8

const (
	STRING CType = iota + 1
	LIST
	DICT
)

func (c *Cache) Get(key string) (ClockedValue, bool) {
	c.mx.Lock()
	defer c.mx.Unlock()
	v, ok := c.store[key]
	return v, ok
}

func (c *Cache) Set(key string, value *CacheValue, context *vector_clock.VC) {
	c.mx.Lock()
	defer c.mx.Unlock()
	c.store[key] = ClockedValue{
		CacheValue: value,
		VC:         context,
	}
}

// TODO: Is that logic should be used for regular set opration?
// Read local value
// If local Version Vector descends incoming Version Vector ignore write (you’ve seen it!)
// If Incoming Version Vector descends local Version Vector overwrite local value with new one
// If Incoming Version Vector is concurrent with local Version Vector, merge values
func (c *Cache) ReplicaSet(key string, value *CacheValue, vc *vector_clock.VC, nodeName string) {
	log.Infof("Incoming ReplicaSet %v %v %v", key, *value, *vc)
	c.mx.Lock()
	defer c.mx.Unlock()
	localValue, ok := c.store[key]
	if ok {
		switch vector_clock.Compare(localValue.VC, vc) {
		case -1:
			log.Infof("Override value on replica node %v", value)
			c.store[key] = ClockedValue{value, vc}
		case 0:
			// merge values and vector clocks then increment result vc
			log.Info("Concurrent write to cache. Trying to merge values...")
			localValue.CacheValue.Merge(value)
			newVc := vector_clock.Merge(localValue.VC, vc)
			newVc.Incr(nodeName)
			c.store[key] = ClockedValue{localValue.CacheValue, newVc}
		}
	} else {
		log.Infof("Set value on replica node %v", value)
		c.store[key] = ClockedValue{value, vc}
	}
}
