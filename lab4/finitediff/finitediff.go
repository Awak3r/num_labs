package finitediff

import (
	"fmt"
	"math"

	"lab4/lu"
	"lab4/matrix"
)

func Solve(p, q, r func(float64) float64, x0, xk, h float64,
	a0, b0, g0, a1, b1, g1 float64) ([]float64, []float64, error) {
	if h <= 0 || xk <= x0 {
		return nil, nil, fmt.Errorf("finitediff: некорректные параметры: x0 = %g, xk = %g, h = %g", x0, xk, h)
	}
	n := int(math.Round((xk - x0) / h))
	if n < 2 {
		return nil, nil, fmt.Errorf("finitediff: нужно не менее 2 шагов, есть %d", n)
	}
	xs := make([]float64, n+1)
	for i := range xs {
		xs[i] = x0 + float64(i)*h
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
		A.Set(i, i-1, 1/(h*h)-p(x)/(2*h))
		A.Set(i, i, -2/(h*h)+q(x))
		A.Set(i, i+1, 1/(h*h)+p(x)/(2*h))
		rhs.Set(i, 0, r(x))
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
	res, err := lu.Solve(A, rhs)
	if err != nil {
		return nil, nil, fmt.Errorf("finitediff: %w", err)
	}
	ys := make([]float64, n+1)
	for i := range ys {
		ys[i] = res.Get(i, 0)
	}
	return xs, ys, nil
}
