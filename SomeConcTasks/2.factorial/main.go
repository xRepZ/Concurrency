package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func main() {
	var wg sync.WaitGroup
	var mu sync.Mutex
	var i, j, k uint64

	wg.Add(2)
	next := true
	go func() {
		defer wg.Done()
		for atomic.LoadUint64(&i) < 1000 {
			atomic.AddUint64(&k, 1)
			mu.Lock()
			if next {
				fmt.Println("Ping")
				atomic.AddUint64(&i, 1)
			}
			next = false
			mu.Unlock()
		}
	}()
	go func() {
		defer wg.Done()

		for atomic.LoadUint64(&i) < 1000 {
			atomic.AddUint64(&j, 1)
			mu.Lock()
			if !next {
				fmt.Println("Pong")
				atomic.AddUint64(&i, 1)
			}
			next = true
			mu.Unlock()

		}

	}()
	wg.Wait()
	fmt.Println(j, k)

}
