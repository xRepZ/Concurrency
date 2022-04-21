package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(3)

	start := time.Now()
	go func() {
		defer wg.Done()
		time.Sleep(time.Second)
	}()

	go func() {
		defer wg.Done()
		time.Sleep(2 * time.Second)
	}()

	go func() {
		defer wg.Done()
		time.Sleep(3 * time.Second)
	}()

	//time.Sleep(d)
	wg.Wait()
	elapsedTime := time.Since(start)

	fmt.Println("Total Time For Execution: " + elapsedTime.String())

}
