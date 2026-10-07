package system

import (
	"fmt"
	"math"

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
		dx, err := luSolve(j, fx)
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

func luSolve(a, b *matrix.Matrix) (*matrix.Matrix, error) {
	n := a.Rows()
	u := a.Copy()
	l := matrix.MakeIdentityMatrix(n)
	for k := 0; k < n; k++ {
		pivot := k
		max := math.Abs(u.Get(k, k))
		for i := k + 1; i < n; i++ {
			if v := math.Abs(u.Get(i, k)); v > max {
				max = v
				pivot = i
			}
		}
		if max < matrix.Eps {
			return nil, fmt.Errorf("вырожденный якобиан (столбец %d)", k)
		}
		if pivot != k {
			u.SwapRows(pivot, k)
			b.SwapRows(pivot, k)
			for j := 0; j < k; j++ {
				tmp := l.Get(pivot, j)
				l.Set(pivot, j, l.Get(k, j))
				l.Set(k, j, tmp)
			}
		}
		for i := k + 1; i < n; i++ {
			l.Set(i, k, u.Get(i, k)/u.Get(k, k))
			for j := k; j < n; j++ {
				u.Set(i, j, u.Get(i, j)-l.Get(i, k)*u.Get(k, j))
			}
		}
	}
	y := matrix.Zeros(n, 1)
	for i := 0; i < n; i++ {
		sum := b.Get(i, 0)
		for j := 0; j < i; j++ {
			sum -= l.Get(i, j) * y.Get(j, 0)
		}
		y.Set(i, 0, sum)
	}
	x := matrix.Zeros(n, 1)
	for i := n - 1; i >= 0; i-- {
		sum := y.Get(i, 0)
		for j := i + 1; j < n; j++ {
			sum -= u.Get(i, j) * x.Get(j, 0)
		}
		x.Set(i, 0, sum/u.Get(i, i))
	}
	return x, nil
}
