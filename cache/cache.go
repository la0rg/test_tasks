package cache

import (
	"sync"

	"github.com/la0rg/test_tasks/vector_clock"
)

type Cache struct {
	mx    sync.Mutex
	store map[string]CacheValue
}

func NewCache() *Cache {
	return &Cache{store: make(map[string]CacheValue)}
}

type CType uint8

const (
	STRING CType = iota + 1
	LIST
	DICT
)

func (c *Cache) Get(key string) (CacheValue, bool) {
	c.mx.Lock()
	defer c.mx.Unlock()
	v, ok := c.store[key]
	return v, ok
}

func (c *Cache) Set(key string, value CacheValue, context *vector_clock.VC) {
	c.mx.Lock()
	defer c.mx.Unlock()
	c.store[key] = value
}
