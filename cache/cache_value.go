package cache

import (
	"errors"

	"github.com/la0rg/test_tasks/vector_clock"
)

// ClockedValue contains value and vector clock associated with it.
type ClockedValue struct {
	*CacheValue
	*vector_clock.VC
}

type CacheValue struct {
	cType       CType
	StringValue string
	ListValue   []CacheValue
	DictValue   map[string]CacheValue
}

var (
	ValueIsNotAccessable error = errors.New("CacheError.ValueIsNotAccessable")
	TypeIsNotSupported   error = errors.New("CacheError.TypeIsNotSupported")
)

func (c *CacheValue) GetType() CType {
	return c.cType
}

func (c *CacheValue) SetString(value string) {
	c.cType = STRING
	c.StringValue = value
	c.ListValue = nil
	c.DictValue = nil
}

func (c *CacheValue) GetString() (string, error) {
	if c.cType != STRING {
		return "", ValueIsNotAccessable
	}
	return c.StringValue, nil
}

func (c *CacheValue) SetList(value []CacheValue) {
	c.cType = LIST
	c.ListValue = value
	c.DictValue = nil
}

func (c *CacheValue) GetList() ([]CacheValue, error) {
	if c.cType != LIST {
		return nil, ValueIsNotAccessable
	}
	return c.ListValue, nil
}

func (c *CacheValue) SetDict(value map[string]CacheValue) {
	c.cType = DICT
	c.DictValue = value
	c.ListValue = nil
}

func (c *CacheValue) GetDict() (map[string]CacheValue, error) {
	if c.cType != DICT {
		return nil, ValueIsNotAccessable
	}
	return c.DictValue, nil
}

func (c *CacheValue) Set(value interface{}) error {
	switch v := value.(type) {
	case string:
		c.SetString(v)
	case []CacheValue:
		c.SetList(v)
	case map[string]CacheValue:
		c.SetDict(v)
	default:
		return TypeIsNotSupported
	}
	return nil
}

func (c *CacheValue) Merge(value *CacheValue) error {
	switch value.cType {
	case STRING:
		if c.cType == STRING {
			c.SetString(c.StringValue + ";merged;" + value.StringValue)
		} else {
			c.SetString(value.StringValue)
		}
	case LIST:
		if c.cType == LIST {
			for _, v := range value.ListValue {
				c.ListValue = append(c.ListValue, v)
			}
		} else {
			c.SetList(value.ListValue)
		}
	case DICT:
		if c.cType == DICT {
			for k, v := range value.DictValue {
				c.DictValue[k] = v
			}
		} else {
			c.SetDict(value.DictValue)
		}
	default:
		return TypeIsNotSupported
	}
	return nil
}
