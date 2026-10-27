package shooting

import (
	"fmt"
	"math"

	"lab4/ode"
)

func Solve(f ode.Func, x0, xk, du0, s0, s1, h, eps float64, maxIter int,
	g func(u []float64) float64) (float64, int, []float64, [][]float64, error) {
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
		xs, us = xx, uu
		return g(uu[len(uu)-1]), nil
	}
	g0, err := G(s0)
	if err != nil {
		return 0, 0, nil, nil, fmt.Errorf("shooting: %w", err)
	}
	g1, err := G(s1)
	if err != nil {
		return 0, 0, nil, nil, fmt.Errorf("shooting: %w", err)
	}

	for iters := 0; iters < maxIter; iters++ {
		if math.Abs(g1-g0) < 1e-300 {
			return 0, 0, nil, nil, fmt.Errorf("shooting: секущие выродились: g(s0) = g(s1) = %g", g1)
		}
		s2 := s1 - g1*(s1-s0)/(g1-g0)
		g2, err := G(s2)
		if err != nil {
			return 0, 0, nil, nil, fmt.Errorf("shooting: %w", err)
		}
		if math.Abs(s2-s1) < eps {
			return s2, iters + 1, xs, us, nil
		}
		s0, g0 = s1, g1
		s1, g1 = s2, g2
	}
	return 0, 0, nil, nil, fmt.Errorf("shooting: секущие не сошлись за %d итераций", maxIter)
}
