package main

import (
	"fmt"
)

func main() {
	a := [10]int32{1, 2, 3, 4, 5, 6, 7, 7, 8, 9}

	b := []int32{4, 12, 6, 3, 11, 6, 9, 10, 3, 9}

	result := make([]int32, 10)

	for i := 0; i < len(a); i++ {
		result[i] = a[i] + b[i]
	}

	fmt.Println("Результат суми:", result)
}
