package demo

import (
	"fmt"
	"math"
)

// rungeError — оценка погрешности решения с шагом h1 по правилу Рунге
// (метод Рунге-Ромберга): ε = max|y_h1 − y_h2| / (1 − 2^(−p)),
// где p — порядок метода, h2 = h1/k (k целое) — более мелкий шаг.
// Сравнение — по общим узлам сеток.
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

// maxAbsDiff — макс. по узлам |y_i − y_точное(x_i)|.
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

// col — компонента y (первая) из массива векторов решения.
func col(us [][]float64) []float64 {
	ys := make([]float64, len(us))
	for i, u := range us {
		ys[i] = u[0]
	}
	return ys
}
