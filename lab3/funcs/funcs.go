package funcs

import "math"

func Tan(x float64) float64 { return math.Tan(x) }

func Integrand(x float64) float64 {
	return x / math.Pow(3*x+4, 3)
}

func ExactIntegral(x0, xk float64) float64 {
	antider := func(x float64) float64 {
		u := 3*x + 4
		return -1/(9*u) + 2/(9*u*u)
	}
	return antider(xk) - antider(x0)
}
