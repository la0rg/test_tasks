package cache

import (
	"sync"

	"github.com/la0rg/test_tasks/vector_clock"
)

type Cache struct {
	mx    sync.Mutex
	store map[string]CValue
}

func NewCache() *Cache {
	return &Cache{store: make(map[string]CValue)}
}

type CType uint8

const (
	STRING CType = iota + 1
	LIST
	DICT
)

type CValue interface {
	GetType() CType

	SetString(value string)
	GetString() string

	SetList(value []CValue)
	GetList() []CValue

	SetDict(value map[CValue]CValue)
	GetDict() map[CValue]CValue
}

func (c *Cache) Get(key string) (CValue, bool) {
	c.mx.Lock()
	defer c.mx.Unlock()
	v, ok := c.store[key]
	return v, ok
}

func (c *Cache) Set(key string, value CValue, context *vector_clock.VC) {
	c.mx.Lock()
	defer c.mx.Unlock()
	c.store[key] = value
}
