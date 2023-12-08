package main

import (
	"log"
	"math/rand"
	"sync"
)

func NewMap() *sync.Map {
	return new(sync.Map)
}

func testMap() {
	syncMap := new(sync.Map)

	syncMap.Store("key", rand.Uint64())

	generalMap := make(map[string]uint64)
	generalMap["key"] = rand.Uint64()
	for i := 0; i < 100000; i++ {
		go func() {
			value, ok := syncMap.Load("key")
			if ok {
				log.Println(value)
				syncMap.CompareAndSwap("key", value, rand.Uint64())
			} else {
				log.Fatal("not found")
			}
		}()
	}

	for i := 0; i < 100000; i++ {
		go func() {
			value, ok := generalMap["key"]
			if ok {
				log.Println(value)
				generalMap["key"] = rand.Uint64()
			} else {
				log.Fatal("not found")
			}
		}()
	}
}
