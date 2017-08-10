package cache

import (
	"errors"
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

// Set reads local value
// If local Version Vector descends incoming Version Vector ignore write (you’ve seen it!)
// If Incoming Version Vector descends local Version Vector overwrite local value with new one
// If Incoming Version Vector is concurrent with local Version Vector returns error
func (c *Cache) Set(key string, value *CacheValue, context *vector_clock.VC) error {
	log.Infof("Set %v %v", key, *value)
	c.mx.Lock()
	defer c.mx.Unlock()
	localValue, ok := c.store[key]

	if ok {
		switch vector_clock.Compare(localValue.VC, context) {
		case 0:
			log.Warn("Conflict write")
			return errors.New("Conflict write")

		case 1:
			log.Info("Ignore write")
			return nil
		}
	}

	log.Infof("Set value on replica node %v", value)
	c.store[key] = ClockedValue{value, context}
	return nil
}
