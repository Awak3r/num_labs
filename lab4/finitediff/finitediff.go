package finitediff

import (
	"fmt"
	"math"

	"lab4/lu"
	"lab4/matrix"
)

func Solve(p, q, r func(float64) float64, x0, xk, h float64,
	a0, b0, g0, a1, b1, g1 float64) ([]float64, []float64, error) {
	if !finite(x0, xk, h, a0, b0, g0, a1, b1, g1) || h <= 0 || xk <= x0 || p == nil || q == nil || r == nil {
		return nil, nil, fmt.Errorf("finitediff: некорректные параметры: x0 = %g, xk = %g, h = %g", x0, xk, h)
	}
	steps := (xk - x0) / h
	rounded := math.Round(steps)
	if !finite(steps) || rounded >= float64(int(^uint(0)>>1)) || math.Abs(steps-rounded) > 1e-9 {
		return nil, nil, fmt.Errorf("finitediff: шаг %g не делит отрезок [%g, %g] или сетка слишком велика", h, x0, xk)
	}
	n := int(rounded)
	if n < 2 {
		return nil, nil, fmt.Errorf("finitediff: нужно не менее 2 шагов, есть %d", n)
	}
	xs := make([]float64, n+1)
	for i := range xs {
		xs[i] = x0 + float64(i)*h
		if i == n {
			xs[i] = xk
		}
		if !finite(xs[i]) || (i > 0 && xs[i] <= xs[i-1]) {
			return nil, nil, fmt.Errorf("finitediff: непредставимый узел сетки %d", i)
		}
	}
	A := matrix.Zeros(n+1, n+1)
	rhs := matrix.Column(make([]float64, n+1))

	if b0 == 0 {
		A.Set(0, 0, a0)
		rhs.Set(0, 0, g0)
	} else {
		A.Set(0, 0, a0-3*b0/(2*h))
		A.Set(0, 1, 2*b0/h)
		A.Set(0, 2, -b0/(2*h))
		rhs.Set(0, 0, g0)
	}

	for i := 1; i < n; i++ {
		x := xs[i]
		px, qx, rx := p(x), q(x), r(x)
		if !finite(px, qx, rx) {
			return nil, nil, fmt.Errorf("finitediff: неконечный коэффициент при x = %g", x)
		}
		A.Set(i, i-1, 1/(h*h)-px/(2*h))
		A.Set(i, i, -2/(h*h)+qx)
		A.Set(i, i+1, 1/(h*h)+px/(2*h))
		rhs.Set(i, 0, rx)
	}

	if b1 == 0 {
		A.Set(n, n, a1)
		rhs.Set(n, 0, g1)
	} else {
		A.Set(n, n, a1+3*b1/(2*h))
		A.Set(n, n-1, -2*b1/h)
		A.Set(n, n-2, b1/(2*h))
		rhs.Set(n, 0, g1)
	}
	for i := range xs {
		if !finite(rhs.Get(i, 0)) {
			return nil, nil, fmt.Errorf("finitediff: неконечная правая часть в строке %d", i)
		}
		for j := range xs {
			if !finite(A.Get(i, j)) {
				return nil, nil, fmt.Errorf("finitediff: неконечный элемент системы [%d][%d]", i, j)
			}
		}
	}
	res, err := lu.Solve(A, rhs)
	if err != nil {
		return nil, nil, fmt.Errorf("finitediff: %w", err)
	}
	ys := make([]float64, n+1)
	for i := range ys {
		ys[i] = res.Get(i, 0)
		if !finite(ys[i]) {
			return nil, nil, fmt.Errorf("finitediff: неконечное решение при x = %g", xs[i])
		}
	}
	return xs, ys, nil
}

func finite(values ...float64) bool {
	for _, v := range values {
		if math.IsNaN(v) || math.IsInf(v, 0) {
			return false
		}
	}
	return true
}
