package main

import "sync"

type Singleton struct {
	name string
}

var once sync.Once
var instance *Singleton

func GetInstance() *Singleton {
	// Do() 메서드로 싱글톤 인스턴스를 생성합니다.
	once.Do(func() {
		instance = &Singleton{
			name: "Singleton",
		}
	})

	return instance
}
