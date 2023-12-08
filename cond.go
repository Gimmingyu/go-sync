package main

import (
	"fmt"
	"sync"
	"time"
)

func testCond() {
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
