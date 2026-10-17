package finitediff

import (
	"fmt"
	"math"

	"lab4/lu"
	"lab4/matrix"
)

// Solve решает линейную краевую задачу
//
//	y'' + p(x)·y' + q(x)·y = r(x),
//	a0·y(x0) + b0·y'(x0) = g0,   a1·y(xk) + b1·y'(xk) = g1
//
// конечно-разностным методом: во внутренних узлах — центральная разность,
// на концах производные аппроксимируются односторонними разностями 2-го
// порядка, поэтому схема имеет 2-й порядок (p = 2). Полученная система
// решается LU-разложением.
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
	// левое условие: a0·y0 + b0·(−3y0 + 4y1 − y2)/(2h) = g0
	if b0 == 0 {
		A.Set(0, 0, a0)
		rhs.Set(0, 0, g0)
	} else {
		A.Set(0, 0, a0-3*b0/(2*h))
		A.Set(0, 1, 2*b0/h)
		A.Set(0, 2, -b0/(2*h))
		rhs.Set(0, 0, g0)
	}
	// внутренние узлы: (y_{i−1} − 2y_i + y_{i+1})/h² + p_i(y_{i+1} − y_{i−1})/(2h) + q_i·y_i = r_i
	for i := 1; i < n; i++ {
		x := xs[i]
		A.Set(i, i-1, 1/(h*h)-p(x)/(2*h))
		A.Set(i, i, -2/(h*h)+q(x))
		A.Set(i, i+1, 1/(h*h)+p(x)/(2*h))
		rhs.Set(i, 0, r(x))
	}
	// правое условие: a1·yN + b1·(3yN − 4y_{N−1} + y_{N−2})/(2h) = g1
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
