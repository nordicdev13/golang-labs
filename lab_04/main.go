package main

import (
	"fmt"
)

func main() {
	numsChan := generate(1, 100)

	evenChan := filterEven(numsChan)

	squaredChan := square(evenChan)

	total := sum(squaredChan)

	fmt.Printf("Cума квадратів парних чисел від 1 до 100: %.0f\n", total)
}

func generate(start, end int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for i := start; i <= end; i++ {
			out <- i
		}
	}()
	return out
}

func filterEven(in <-chan int) <-chan int {
	out := make(chan int, 10)
	go func() {
		defer close(out)
		for n := range in {
			if n%2 == 0 {
				out <- n
			}
		}
	}()
	return out
}

func square(in <-chan int) <-chan float64 {
	out := make(chan float64)
	go func() {
		defer close(out)
		for n := range in {
			out <- float64(n * n)
		}
	}()
	return out
}

func sum(in <-chan float64) float64 {
	var total float64
	for n := range in {
		total += n
	}
	return total
}
