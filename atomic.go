package main

import (
	"sync/atomic"
)

type safeCache struct {
	data atomic.Value // Stores the actual map data
}

func newSafeCache() *safeCache {
	sc := &safeCache{}
	sc.data.Store(make(map[string]string))
	return sc
}

func (sc *safeCache) set(key string, value string) {
	// Load current value of the cache
	data := sc.data.Load().(map[string]string)

	// Create a new map with existing data and the new key-value pair
	newData := make(map[string]string)
	for k, v := range data {
		newData[k] = v
	}
	newData[key] = value

	// Atomically replace the current map with the new map
	sc.data.Swap(newData)
}

func (sc *safeCache) get(key string) (string, bool) {
	data := sc.data.Load().(map[string]string)
	value, ok := data[key]
	return value, ok
}
