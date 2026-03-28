package calc

import (
	"errors"
)

// Завдання 1

func Sum(nums ...float64) float64 {
	total := 0.0
	for _, n := range nums {
		total += n
	}
	return total
}

func Max(nums ...float64) float64 {
	if len(nums) == 0 {
		return 0
	}
	res := nums[0]
	for _, n := range nums {
		if n > res {
			res = n
		}
	}
	return res
}

func Min(nums ...float64) float64 {
	if len(nums) == 0 {
		return 0
	}
	res := nums[0]
	for _, n := range nums {
		if n < res {
			res = n
		}
	}
	return res
}

func Divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("ділення на нуль неможливе")
	}
	return a / b, nil
}

// Завдання 2

type Calculator interface {
	Sum(nums ...float64) float64
	Max(nums ...float64) float64
	Min(nums ...float64) float64
	Divide(a, b float64) (float64, error)
}

type Calc struct{}

func NewCalc() *Calc {
	return &Calc{}
}

func (c *Calc) Sum(nums ...float64) float64 {
	total := 0.0
	for _, n := range nums {
		total += n
	}
	return total
}

func (c *Calc) Max(nums ...float64) float64 {
	if len(nums) == 0 {
		return 0
	}
	res := nums[0]
	for _, n := range nums {
		if n > res {
			res = n
		}
	}
	return res
}

func (c *Calc) Min(nums ...float64) float64 {
	if len(nums) == 0 {
		return 0
	}
	res := nums[0]
	for _, n := range nums {
		if n < res {
			res = n
		}
	}
	return res
}

func (c *Calc) Divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("ділення на нуль")
	}
	return a / b, nil
}
