package cache

import (
	"sync"

	"github.com/la0rg/test_tasks/vector_clock"
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
