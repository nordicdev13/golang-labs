package main

import (
	"fmt"
	"math"
)

type Shape interface {
	Area() float64
	Perimeter() float64
}

type Circle struct {
	Radius float64
}

func (c Circle) Area() float64 {
	return math.Pi * math.Pow(c.Radius, 2)
}
func (c Circle) Perimeter() float64 {
	return 2 * math.Pi * c.Radius
}

type Rectangle struct {
	Width, Height float64
}

func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

type Triangle struct {
	A, B, C float64
}

func (t Triangle) Area() float64 {
	p := t.Perimeter() / 2
	return math.Sqrt(p * (p - t.A) * (p - t.B) * (p - t.C))
}

func (t Triangle) Perimeter() float64 {
	return t.A + t.B + t.C
}

func main() {
	c := Circle{Radius: 5}
	r := Rectangle{Width: 10, Height: 4}
	t := Triangle{A: 3, B: 4, C: 5}

	shapes := []Shape{c, r, t}

	fmt.Println("=== Результати обчислень ===")
	for _, s := range shapes {
		fmt.Printf("Фігура: %T\n", s)
		fmt.Printf("Площа: %.2f\n", s.Area())
		fmt.Printf("Периметр: %.2f\n\n", s.Perimeter())
	}
}
