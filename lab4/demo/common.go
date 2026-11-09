package demo

import (
	"fmt"
	"math"
)

func rungeError(h1, h2 float64, p int, xs1, ys1, xs2, ys2 []float64) (float64, error) {
	if !finite(h1, h2) || h1 <= 0 || h2 <= 0 || p <= 0 {
		return 0, fmt.Errorf("Рунге-Ромберг: нужны конечные положительные шаги и p > 0")
	}
	if len(xs1) < 2 || len(xs2) < 2 || len(xs1) != len(ys1) || len(xs2) != len(ys2) {
		return 0, fmt.Errorf("Рунге-Ромберг: неверные размеры сеток и решений")
	}
	k := (len(xs2) - 1) / (len(xs1) - 1)
	if k <= 1 || (len(xs2)-1)%(len(xs1)-1) != 0 || !finite(h1/h2) || math.Abs(h1/h2-float64(k)) > 1e-9 {
		return 0, fmt.Errorf("Рунге-Ромберг: сетки должны покрывать один отрезок, h1/h2 — целое k > 1: h1 = %g, h2 = %g", h1, h2)
	}
	for _, grid := range []struct {
		h      float64
		xs, ys []float64
	}{{h1, xs1, ys1}, {h2, xs2, ys2}} {
		for i, x := range grid.xs {
			if !finite(x, grid.ys[i]) || (i > 0 && x <= grid.xs[i-1]) || !(math.Abs((x-grid.xs[0])/grid.h-float64(i)) <= 1e-9) {
				return 0, fmt.Errorf("Рунге-Ромберг: некорректная сетка или решение в узле %d", i)
			}
		}
	}
	max := 0.0
	for i := range xs1 {
		j := i * k
		if !(math.Abs((xs1[i]-xs2[j])/h2) <= 1e-9) {
			return 0, fmt.Errorf("Рунге-Ромберг: узлы сеток не совпадают: %g и %g", xs1[i], xs2[j])
		}
		d := math.Abs(ys1[i] - ys2[j])
		if d > max {
			max = d
		}
	}
	err := max / (1 - math.Pow(float64(k), -float64(p)))
	if !finite(err) {
		return 0, fmt.Errorf("Рунге-Ромберг: неконечная оценка погрешности")
	}
	return err, nil
}

func finite(values ...float64) bool {
	for _, v := range values {
		if math.IsNaN(v) || math.IsInf(v, 0) {
			return false
		}
	}
	return true
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
