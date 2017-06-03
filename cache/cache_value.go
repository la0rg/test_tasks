package cache

import "errors"

type CacheValue struct {
	cType       CType
	stringValue string
	listValue   []CValue
	dictValue   map[CValue]CValue
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
	c.stringValue = value
}

func (c *CacheValue) GetString() (string, error) {
	if c.cType != STRING {
		return "", ValueIsNotAccessable
	}
	return c.stringValue, nil
}

func (c *CacheValue) SetList(value []CValue) {
	c.cType = LIST
	c.listValue = value
}

func (c *CacheValue) GetList() ([]CValue, error) {
	if c.cType != LIST {
		return nil, ValueIsNotAccessable
	}
	return c.listValue, nil
}

func (c *CacheValue) SetDict(value map[CValue]CValue) {
	c.cType = DICT
	c.dictValue = value
}

func (c *CacheValue) GetDict() (map[CValue]CValue, error) {
	if c.cType != DICT {
		return nil, ValueIsNotAccessable
	}
	return c.dictValue, nil
}

func (c *CacheValue) Set(value interface{}) error {
	switch v := value.(type) {
	case string:
		c.SetString(v)
	case []CValue:
		c.SetList(v)
	case map[CValue]CValue:
		c.SetDict(v)
	default:
		return TypeIsNotSupported
	}
	return nil
}
