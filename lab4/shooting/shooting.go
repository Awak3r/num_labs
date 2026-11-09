package shooting

import (
	"fmt"
	"math"

	"lab4/ode"
)

func Solve(f ode.Func, x0, xk, du0, s0, s1, h, eps float64, maxIter int,
	g func(u []float64) float64) (float64, int, []float64, [][]float64, error) {
	if !finite(x0, xk, du0, s0, s1, h, eps) || eps <= 0 || maxIter <= 0 || f == nil || g == nil {
		return 0, 0, nil, nil, fmt.Errorf("shooting: нужны конечные параметры, eps > 0, maxIter > 0 и заданные функции")
	}
	if s0 == s1 {
		return 0, 0, nil, nil, fmt.Errorf("shooting: начальные приближения s0 и s1 совпадают")
	}
	if h <= 0 || xk <= x0 {
		return 0, 0, nil, nil, fmt.Errorf("shooting: некорректные параметры: x0 = %g, xk = %g, h = %g", x0, xk, h)
	}

	var xs []float64
	var us [][]float64
	G := func(s float64) (float64, error) {
		xx, uu, err := ode.RK4(f, x0, xk, h, []float64{s, du0})
		if err != nil {
			return 0, err
		}
		residual := g(uu[len(uu)-1])
		if !finite(residual) {
			return 0, fmt.Errorf("неконечная невязка при s = %g", s)
		}
		xs, us = xx, uu
		return residual, nil
	}
	g0, err := G(s0)
	if err != nil {
		return 0, 0, nil, nil, fmt.Errorf("shooting: %w", err)
	}
	if g0 == 0 {
		return s0, 0, xs, us, nil
	}
	g1, err := G(s1)
	if err != nil {
		return 0, 0, nil, nil, fmt.Errorf("shooting: %w", err)
	}
	if g1 == 0 || (math.Abs(g1) <= eps && math.Abs(s1-s0) < eps) {
		return s1, 0, xs, us, nil
	}

	for iters := 0; iters < maxIter; iters++ {
		if !finite(g1-g0, s1-s0) {
			return 0, 0, nil, nil, fmt.Errorf("shooting: переполнение разностей секущих")
		}
		if math.Abs(g1-g0) < 1e-300 {
			if math.Abs(g1) <= eps {
				return s1, iters, xs, us, nil
			}
			return 0, 0, nil, nil, fmt.Errorf("shooting: секущие выродились: g(s0) = g(s1) = %g", g1)
		}
		s2 := s1 - g1*(s1-s0)/(g1-g0)
		if !finite(s2) {
			return 0, 0, nil, nil, fmt.Errorf("shooting: неконечный параметр секущих")
		}
		if s2 == s1 {
			if math.Abs(g1) <= eps {
				return s1, iters + 1, xs, us, nil
			}
			return 0, 0, nil, nil, fmt.Errorf("shooting: стагнация секущих при s = %g, невязка = %g", s1, g1)
		}
		g2, err := G(s2)
		if err != nil {
			return 0, 0, nil, nil, fmt.Errorf("shooting: %w", err)
		}
		// Малое изменение параметра само по себе не гарантирует выполнения краевого условия.
		if math.Abs(g2) <= eps && (math.Abs(s2-s1) < eps || g2 == 0) {
			return s2, iters + 1, xs, us, nil
		}
		s0, g0 = s1, g1
		s1, g1 = s2, g2
	}
	return 0, 0, nil, nil, fmt.Errorf("shooting: секущие не сошлись за %d итераций, невязка = %g", maxIter, g1)
}

func finite(values ...float64) bool {
	for _, v := range values {
		if math.IsNaN(v) || math.IsInf(v, 0) {
			return false
		}
	}
	return true
}
