package demo

import (
	"fmt"
	"math"
)

func rungeError(h1, h2 float64, p int, xs1, ys1, xs2, ys2 []float64) (float64, error) {
	k := int(math.Round(h1 / h2))
	if k < 1 || math.Abs(float64(k)*h2-h1) > 1e-9 {
		return 0, fmt.Errorf("шаги не кратны: h1 = %g, h2 = %g", h1, h2)
	}
	max := 0.0
	for i := range xs1 {
		j := i * k
		if j >= len(ys2) {
			break
		}
		d := math.Abs(ys1[i] - ys2[j])
		if d > max {
			max = d
		}
	}
	return max / (1 - math.Pow(2, float64(-p))), nil
}

func maxAbsDiff(xs, ys []float64, exact func(float64) float64) float64 {
	max := 0.0
	for i := range xs {
		d := math.Abs(ys[i] - exact(xs[i]))
		if d > max {
			max = d
		}
	}
	return max
}

func col(us [][]float64) []float64 {
	ys := make([]float64, len(us))
	for i, u := range us {
		ys[i] = u[0]
	}
	return ys
}
