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
