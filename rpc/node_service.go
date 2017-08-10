package rpc

import (
	"github.com/la0rg/test_tasks/cache"
)

// Go transforms proto struct to original go struct
func (m *ClockedValue) Go() *cache.ClockedValue {
	return &cache.ClockedValue{
		CacheValue: m.GetValue().Go(),
		VC:         m.GetVectorClock().Go(),
	}
}

// Go transforms proto struct to original go struct
func (value *CacheValue) Go() *cache.CacheValue {
	cValue := &cache.CacheValue{}
	switch cache.CType(value.CType) {
	case cache.STRING:
		cValue.SetString(value.StringValue)
	case cache.LIST:
		list := make([]cache.CacheValue, len(value.ListValue))
		for i, v := range value.ListValue {
			list[i] = *(v.Go())
		}
		cValue.SetList(list)
	case cache.DICT:
		dict := make(map[string]cache.CacheValue, len(value.DictValue))
		for k, v := range value.DictValue {
			dict[k] = *(v.Go())
		}
		cValue.SetDict(dict)
	}
	return cValue
}

type ProtoClockedValue cache.ClockedValue

func (c ProtoClockedValue) Proto() *ClockedValue {
	return &ClockedValue{
		Value:       ProtoCacheValue(*c.CacheValue).Proto(),
		VectorClock: &VC{c.VC.GetStore()},
	}
}

type ProtoCacheValue cache.CacheValue

func (c ProtoCacheValue) Proto() *CacheValue {
	cValue := &CacheValue{}
	switch c.CType {
	case cache.STRING:
		cValue.StringValue = c.StringValue
		cValue.CType = uint32(cache.STRING)
	case cache.LIST:
		list := make([]*CacheValue, len(c.ListValue))
		for i, v := range c.ListValue {
			list[i] = ProtoCacheValue(v).Proto()
		}
		cValue.ListValue = list
		cValue.CType = uint32(cache.LIST)
	case cache.DICT:
		dict := make(map[string]*CacheValue, len(c.DictValue))
		for k, v := range c.DictValue {
			dict[k] = ProtoCacheValue(v).Proto()
		}
		cValue.DictValue = dict
		cValue.CType = uint32(cache.DICT)
	}
	return cValue
}
