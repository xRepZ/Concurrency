package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	mu := &sync.Mutex{}
	ch := make(chan int)
	sl := make([]int, 0)
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 10; i++ {
			ch <- i
		}
		close(ch)

	}()

	go func() {
		defer wg.Done()
		for n := range ch {
			mu.Lock()
			sl = append(sl, n)
			fmt.Println(sl)
			mu.Unlock()

		}

	}()

	//time.Sleep(d)
	wg.Wait()
	fmt.Println(sl)
}
