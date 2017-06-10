package util

import (
	"encoding/json"
	"fmt"

	"github.com/la0rg/test_tasks/cache"
	"github.com/pkg/errors"
)

type SimpleCacheValue map[string]interface{}

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
