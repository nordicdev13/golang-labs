package main

import (
	"fmt"
	"go-labs/calc"
)

func main() {
	// Завдання 1

	a, b := 10.5, 2.5
	numbers := []float64{1.2, 5.5, 3.8, 10.1, 0.4}

	fmt.Printf("Сума: %.2f\n", calc.Sum(numbers...))

	fmt.Printf("Максимум: %.2f\n", calc.Max(numbers...))
	fmt.Printf("Мінімум: %.2f\n", calc.Min(numbers...))

	res, err := calc.Divide(a, b)
	if err != nil {
		fmt.Println("Помилка:", err)
	} else {
		fmt.Printf("Результат ділення %.1f / %.1f = %.2f\n", a, b, res)
	}

	// Завдання 2

	myCalc := calc.NewCalc()
	runCalculator(myCalc)
}

func runCalculator(c calc.Calculator) {

	vals := []float64{10, 20, 30, -5}

	fmt.Printf("Сума: %.1f\n", c.Sum(vals...))
	fmt.Printf("Максимум: %.1f\n", c.Max(vals...))
	fmt.Printf("Мінімум: %.1f\n", c.Min(vals...))

	res, err := c.Divide(100, 4)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Printf("Результат ділення (100/4): %.1f\n", res)
	}
}
