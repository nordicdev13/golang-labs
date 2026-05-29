package main

import (
	"fmt"
	"sync"
)

func main() {
	evenChan := make(chan int)
	oddChan := make(chan int)

	var counter int
	var mu sync.Mutex

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		for i := 1; i <= 1000; i++ {
			select {
			case val := <-evenChan:
				if val%3 == 0 {
					mu.Lock()
					counter++
					mu.Unlock()
				}
			case val := <-oddChan:
				if val%33 == 0 {
					mu.Lock()
					counter--
					mu.Unlock()
				}
			}
		}
	}()

	for i := 1; i <= 1000; i++ {
		if i%2 == 0 {
			evenChan <- i
		} else {
			oddChan <- i
		}
	}

	wg.Wait()

	fmt.Printf("Фінальне значення counter: %d\n", counter)
}
