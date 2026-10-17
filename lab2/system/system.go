package system

import (
	"fmt"
	"math"

	"lab2/lu"
	"lab2/matrix"
)

type Step struct {
	Iter int
	X    *matrix.Matrix
	Fx   float64
	Diff float64
}

func SimpleIterations(phi func(*matrix.Matrix) *matrix.Matrix, x0 *matrix.Matrix, eps float64, maxIter int) (*matrix.Matrix, int, []Step, error) {
	if phi0 := phi(x0); phi0 == nil {
		return nil, 0, nil, fmt.Errorf("system: phi(x0) вне области определения")
	}
	x := x0.Copy()
	var steps []Step
	for k := 0; k < maxIter; k++ {
		xNew := phi(x)
		if xNew == nil {
			return nil, 0, nil, fmt.Errorf("system: итерация %d: phi(x) вне области определения", k+1)
		}
		diff, err := xNew.Sub(x)
		if err != nil {
			return nil, 0, nil, err
		}
		steps = append(steps, Step{Iter: k + 1, X: xNew.Copy(), Fx: math.NaN(), Diff: diff.MaxAbs()})
		if diff.MaxAbs() < eps {
			return xNew, k + 1, steps, nil
		}
		x = xNew
	}
	return nil, 0, nil, fmt.Errorf("system: простые итерации не сошлись за %d итераций", maxIter)
}

func Newton(f func(*matrix.Matrix) *matrix.Matrix, jac func(*matrix.Matrix) *matrix.Matrix, x0 *matrix.Matrix, eps float64, maxIter int) (*matrix.Matrix, int, []Step, error) {
	x := x0.Copy()
	var steps []Step
	for k := 0; k < maxIter; k++ {
		fx := f(x)
		if fx == nil {
			return nil, 0, nil, fmt.Errorf("system: ньютон: F(x) вне области определения при x = %s", x)
		}
		j := jac(x)
		if j == nil {
			return nil, 0, nil, fmt.Errorf("system: ньютон: якобиан не определён при x = %s", x)
		}
		dx, err := lu.Solve(j, fx)
		if err != nil {
			return nil, 0, nil, fmt.Errorf("system: ньютон: %w", err)
		}
		diff := dx.MaxAbs()
		nx, err := x.Sub(dx)
		if err != nil {
			return nil, 0, nil, err
		}
		if math.IsNaN(nx.Get(0, 0)) || math.IsInf(nx.Get(0, 0), 0) {
			return nil, 0, nil, fmt.Errorf("system: ньютон: итерация %d вышла из области определения", k+1)
		}
		steps = append(steps, Step{Iter: k + 1, X: nx.Copy(), Fx: fx.MaxAbs(), Diff: diff})
		if diff < eps {
			return nx, k + 1, steps, nil
		}
		x = nx
	}
	return nil, 0, nil, fmt.Errorf("system: ньютон не сошёлся за %d итераций", maxIter)
}
