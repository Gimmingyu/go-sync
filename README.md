# go-sync

"sync" 패키지 잘쓰기.

atomic, cond, map, once, (Waitgroup, Mutex, Pool, RWMutex 는 생략)

## atomic 

```go
package main

import (
	"fmt"
	"log"
	"time"
)

func main() {

	cache := newSafeCache()

	for i := 0; i < 5; i++ {
		// Simulating concurrent writes
		for i := 0; i < 5; i++ {
			go func(i int) {
				key := fmt.Sprintf("key%d", i)
				value := fmt.Sprintf("value%d", i)
				cache.set(key, value)
				log.Printf("Set %s = %s\n", key, value)
			}(i)
		}

		time.Sleep(time.Second)

		// Simulating concurrent reads
		for i := 0; i < 10; i++ {
			go func(i int) {
				key := fmt.Sprintf("key%d", i/2)
				value, ok := cache.get(key)
				if ok {
					log.Printf("Get %s = %s\n", key, value)
				} else {
					log.Printf("Get %s = not found\n", key)
				}
			}(i)
		}

		// Wait for a while to observe cache operations
		time.Sleep(2 * time.Second)
	}
}

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

```

## cond

sync.Cond 는 조건 변수를 제공해준다. 

goroutine 이 특정 조건이 충족될 때까지 대기하고, 다른 goroutine 이 특정 조건을 충족시키면 대기하고 있는 goroutine 을 깨워준다.

아래 예시와 같은 pub-sub 패턴을 구현할 때 유용하게 사용할 수 있다.

```go
package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var mu sync.Mutex
	cond := sync.NewCond(&mu)

	queue := make([]int, 0)

	// Producer
	for i := 0; i < 3; i++ {
		go func(i int) {
			for {
				mu.Lock()
				num := time.Now().Nanosecond()
				queue = append(queue, num)
				fmt.Printf("Producer %d produced: %d\n", i, num)
				mu.Unlock()
				cond.Signal()
				time.Sleep(time.Second)
			}
		}(i)
	}

	// Consumer
	for i := 0; i < 3; i++ {
		go func(i int) {
			for {
				mu.Lock()
				for len(queue) == 0 {
					cond.Wait()
				}
				num := queue[0]
				queue = queue[1:]
				fmt.Printf("Consumer %d consumed: %d\n", i, num)
				mu.Unlock()
			}
		}(i)
	}

	// Wait for a while to observe producer and consumer activities
	time.Sleep(10 * time.Second)
}
```

## map

sync.Map 은 언제 쓰면 좋을까? 

일반적으로 map[K]V 의 형태로 map 구조체를 많이 선언하곤 하는데, goroutine 에서 동시에 접근하면 concurrent map read and map write 가 발생한다. 

이를 해결하기 위해 sync.Map 을 사용할 수 있다.

물론 mutex 를 사용해도 되지만, 애초에 thread-safe 하게 설계된 sync.Map 을 사용하는 것이 더 좋다고 본다.

sync.Map 은 map[K]V 의 형태로 사용할 수 없다. 대신, Load, Store, Delete, LoadOrStore, Range 메서드를 사용한다.

```go
package main

import (
    "fmt"
    "sync"
)

func main() {
    var m sync.Map

    m.Store("key", "value")
    v, ok := m.Load("key")
    fmt.Println(v, ok)

    m.Delete("key")
    v, ok = m.Load("key")
    fmt.Println(v, ok)

    v, ok = m.LoadOrStore("key", "value")
    fmt.Println(v, ok)

    m.Range(func(k, v interface{}) bool {
        fmt.Println(k, v)
        return true
    })
}
```
## once

정확히 하나의 작업을 수행하는데 사용한다.

한 번 사용한 후에, 복사해서는 안된다. 아래는 설명 첨부.

> Once is an object that will perform exactly one action.
A Once must not be copied after first use.
In the terminology of the Go memory model, the return from f “synchronizes before” the return from any call of once.Do(f).

즉 하나의 Once 객체에 대해서 여러 번 호출되더라도 첫 번째 호출만 실행된다. 

초기화나 커넥션 연결, 싱글톤 패턴 구현 등 정확히 한 번 실행해야 할 때 유용할 수 있을 것이다. 

```go
package main

import (
    "fmt"
    "sync"
)

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

func main() {
    // 싱글톤 인스턴스를 가져옵니다.
    instance := GetInstance()

    // 싱글톤 인스턴스의 이름을 출력합니다.
    fmt.Println(instance.name)
}
```