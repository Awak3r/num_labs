package equation

import (
	"fmt"
	"math"
)

type Step struct {
	Iter int
	X    float64
	Fx   float64
	Diff float64
}

func domain(x float64) error {
	if math.IsNaN(x) || math.IsInf(x, 0) || math.Abs(x) > 1 {
		return fmt.Errorf("equation: x = %g вне области определения |x| <= 1", x)
	}
	return nil
}

func SimpleIterations(phi func(float64) float64, f func(float64) float64, x0, eps float64, maxIter int) (float64, int, []Step, error) {
	if err := domain(x0); err != nil {
		return 0, 0, nil, fmt.Errorf("equation: x0: %w", err)
	}
	x := x0
	var steps []Step
	for k := 0; k < maxIter; k++ {
		xNew := phi(x)
		if err := domain(xNew); err != nil {
			return 0, 0, nil, fmt.Errorf("equation: итерация %d: %w", k+1, err)
		}
		diff := math.Abs(xNew - x)
		steps = append(steps, Step{Iter: k + 1, X: xNew, Fx: f(xNew), Diff: diff})
		if diff < eps {
			return xNew, k + 1, steps, nil
		}
		x = xNew
	}
	return 0, 0, nil, fmt.Errorf("equation: простые итерации не сошлись за %d итераций", maxIter)
}

func Newton(f func(float64) float64, df func(float64) float64, x0, eps float64, maxIter int) (float64, int, []Step, error) {
	if err := domain(x0); err != nil {
		return 0, 0, nil, fmt.Errorf("equation: x0: %w", err)
	}
	x := x0
	var steps []Step
	for k := 0; k < maxIter; k++ {
		fx := f(x)
		d := df(x)
		if math.Abs(d) < 1e-12 {
			return 0, 0, nil, fmt.Errorf("equation: ньютон: нулевая производная при x = %g", x)
		}
		xNew := x - fx/d
		diff := math.Abs(xNew - x)
		if err := domain(xNew); err != nil {
			return 0, 0, nil, fmt.Errorf("equation: ньютон: итерация %d: %w", k+1, err)
		}
		steps = append(steps, Step{Iter: k + 1, X: xNew, Fx: f(xNew), Diff: diff})
		if diff < eps {
			return xNew, k + 1, steps, nil
		}
		x = xNew
	}
	return 0, 0, nil, fmt.Errorf("equation: ньютон не сошёлся за %d итераций", maxIter)
}
