package main

import (
	"fmt"
	"sync"
)

func main() {
	nameByBook := make(map[string]string)
	names := []string{"Name1", "Name2", "Name3", "Name4", "Name5", "Name6", "Name7"}
	books := []string{"book1", "book2", "book3", "book4", "book5", "book6", "book7"}
	wg := &sync.WaitGroup{}
	mu := &sync.Mutex{}
	wg.Add(4)
	chNames := make(chan string)
	chBooks := make(chan string)
	go func(names, books []string) {
		defer wg.Done()
		for i := 0; i < len(names); i++ {
			chNames <- names[i]
			chBooks <- books[i]
		}
		close(chNames)
		close(chBooks)
	}(names, books)
	go func(nameByBook map[string]string) {
		defer wg.Done()

		for n := range chNames {
			mu.Lock()
			nameByBook[n] = <-chBooks
			mu.Unlock()
		}

	}(nameByBook)
	go func(nameByBook map[string]string) {
		defer wg.Done()
		for n := range chNames {
			mu.Lock()
			nameByBook[n] = <-chBooks
			mu.Unlock()
		}

	}(nameByBook)
	go func(nameByBook map[string]string) {
		defer wg.Done()
		for n := range chNames {
			mu.Lock()
			nameByBook[n] = <-chBooks
			mu.Unlock()
		}

	}(nameByBook)

	wg.Wait()
	for k, v := range nameByBook {
		fmt.Printf(" %v : %v ", k, v)
	}

}
