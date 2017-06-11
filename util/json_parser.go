package util

import (
	"encoding/json"
	"fmt"

	"github.com/la0rg/test_tasks/cache"
	"github.com/pkg/errors"
)

type SimpleCacheValue map[string]interface{}

func CacheValueToJson(value *cache.CacheValue) (interface{}, error) {
	switch value.GetType() {
	case cache.STRING:
		valueString, err := value.GetString()
		if err == cache.ValueIsNotAccessable {
			return nil, errors.New("CacheValue has invalid structure")
		}
		return valueString, nil
	case cache.LIST:
		valueList, err := value.GetList()
		if err == cache.ValueIsNotAccessable {
			return nil, errors.New("CacheValue has invalid structure")
		}
		list := make([]interface{}, len(valueList))
		for i, v := range valueList {
			iterValue, err := CacheValueToJson(&v)
			if err != nil {
				return nil, err
			}
			list[i] = iterValue
		}
		return list, nil
	case cache.DICT:
		valueDict, err := value.GetDict()
		if err == cache.ValueIsNotAccessable {
			return nil, errors.New("CacheValue has invalid structure")
		}
		dict := make(map[string]interface{}, len(valueDict))
		for k, v := range valueDict {
			iterValue, err := CacheValueToJson(&v)
			if err != nil {
				return nil, err
			}
			dict[k] = iterValue
		}
		return dict, nil
	default:
		return nil, errors.New("CacheValue with undefined type")
	}
}

func ParseJson(source []byte) (string, *cache.CacheValue, error) {
	var simpleCacheValue SimpleCacheValue
	var cacheValue cache.CacheValue

	err := json.Unmarshal(source, &simpleCacheValue)
	if err != nil {
		return "", &cacheValue, errors.Wrap(err, "Unable to parse json into CacheValue object")
	}
	key, keyOk := simpleCacheValue["key"]
	var keyString string
	if keyOk {
		keyString, keyOk = key.(string)
	}
	value, valueOk := simpleCacheValue["value"]
	if !keyOk || !valueOk {
		return "", &cacheValue, errors.New("Json object should contain key and value fields")
	}

	cacheValue, err = iterJson(value)
	if err != nil {
		return "", &cacheValue, err
	}

	return keyString, &cacheValue, nil
}

func iterJson(v interface{}) (cache.CacheValue, error) {
	var cacheValue = cache.CacheValue{}
	switch typedValue := v.(type) {
	case string:
		cacheValue.SetString(typedValue)
	case map[string]interface{}:
		cacheValueMap := make(map[string]cache.CacheValue, len(typedValue))
		for k, v := range typedValue {
			value, err := iterJson(v)
			if err != nil {
				return cacheValue, err
			}
			cacheValueMap[k] = value
		}
		cacheValue.SetDict(cacheValueMap)
	case []interface{}:
		cacheValueList := make([]cache.CacheValue, len(typedValue))
		for i, v := range typedValue {
			value, err := iterJson(v)
			if err != nil {
				return cacheValue, err
			}
			cacheValueList[i] = value
		}
		cacheValue.SetList(cacheValueList)
	default:
		return cacheValue, fmt.Errorf("Unable to parse json value: %v", typedValue)
	}
	return cacheValue, nil
}
