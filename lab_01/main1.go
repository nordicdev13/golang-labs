package main

import (
	"fmt"
	"math/rand"
)

func main() {
	var y int32
	x := rand.Int31n(1000)

	switch {
	case x > 30 && x < 60:
		y = (x*x + x) - 100
	default:
		y = 2 * x * x
	}

	fmt.Printf("x: %d\n", x)
	fmt.Printf("Result: %d\n", y)
}
