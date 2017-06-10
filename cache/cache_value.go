package cache

import "errors"

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
