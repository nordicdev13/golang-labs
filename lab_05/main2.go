package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func main() {
	evenChan := make(chan int)
	oddChan := make(chan int)

	var counter atomic.Int64

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		for i := 1; i <= 1000; i++ {
			select {
			case val := <-evenChan:
				if val%3 == 0 {
					counter.Add(1)
				}
			case val := <-oddChan:
				if val%33 == 0 {
					counter.Add(-1)
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

	fmt.Printf("Фінальне значення counter (атомарно): %d\n", counter.Load())
}
