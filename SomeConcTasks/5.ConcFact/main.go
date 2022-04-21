package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	var mu sync.Mutex
	num := 15
	fact := 1
	wg.Add(2)
	next := true
	go func() {
		defer wg.Done()

		for num-1 > 0 {
			mu.Lock()
			if next {
				num--
				if num == 0 {
					return
				}
				fact *= num * (num + 1)
				fmt.Println("Ping")
				num--
				fmt.Println(num)

			}

			next = false
			mu.Unlock()
		}

	}()
	go func() {
		defer wg.Done()

		for num-1 > 0 {
			mu.Lock()

			if !next {
				num--
				if num == 0 {
					return
				}
				fact *= num * (num + 1)
				fmt.Println("Pong")
				num--
				fmt.Println(num)

			}

			next = true
			mu.Unlock()

		}

	}()
	wg.Wait()
	fmt.Println(fact)

}
