package ode

import (
	"fmt"
	"math"
)

type Func func(x float64, u []float64) []float64

func Euler(f Func, x0, xk, h float64, u0 []float64) ([]float64, [][]float64, error) {
	if f == nil || len(u0) == 0 || !finite(u0) {
		return nil, nil, fmt.Errorf("ode: некорректная функция или начальное состояние")
	}
	xs, n, err := nodes(x0, xk, h)
	if err != nil {
		return nil, nil, err
	}
	us := make([][]float64, n+1)
	u := clone(u0)
	us[0] = u
	for k := 0; k < n; k++ {
		du, err := derivative(f, xs[k], u)
		if err != nil {
			return nil, nil, err
		}
		u = axpy(u, du, h)
		if !finite(u) {
			return nil, nil, fmt.Errorf("ode: эйлер: решение расходится при x = %g", xs[k+1])
		}
		us[k+1] = clone(u)
	}
	return xs, us, nil
}

func RK4(f Func, x0, xk, h float64, u0 []float64) ([]float64, [][]float64, error) {
	if f == nil || len(u0) == 0 || !finite(u0) {
		return nil, nil, fmt.Errorf("ode: некорректная функция или начальное состояние")
	}
	xs, n, err := nodes(x0, xk, h)
	if err != nil {
		return nil, nil, err
	}
	us := make([][]float64, n+1)
	u := clone(u0)
	us[0] = u
	for k := 0; k < n; k++ {
		u, err = rk4Step(f, xs[k], u, h)
		if err != nil {
			return nil, nil, err
		}
		us[k+1] = clone(u)
	}
	return xs, us, nil
}

func rk4Step(f Func, x float64, u []float64, h float64) ([]float64, error) {
	k1, err := derivative(f, x, u)
	if err != nil {
		return nil, err
	}
	k2, err := derivative(f, x+h/2, axpy(u, k1, h/2))
	if err != nil {
		return nil, err
	}
	k3, err := derivative(f, x+h/2, axpy(u, k2, h/2))
	if err != nil {
		return nil, err
	}
	k4, err := derivative(f, x+h, axpy(u, k3, h))
	if err != nil {
		return nil, err
	}
	sum := axpy(axpy(axpy(k1, k2, 2), k3, 2), k4, 1)
	u = axpy(u, sum, h/6)
	if !finite(u) {
		return nil, fmt.Errorf("ode: рк4: решение расходится при x = %g", x+h)
	}
	return u, nil
}

func Adams4(f Func, x0, xk, h float64, u0 []float64) ([]float64, [][]float64, error) {
	if f == nil || len(u0) == 0 || !finite(u0) {
		return nil, nil, fmt.Errorf("ode: некорректная функция или начальное состояние")
	}
	xs, n, err := nodes(x0, xk, h)
	if err != nil {
		return nil, nil, err
	}
	if n < 4 {
		return nil, nil, fmt.Errorf("ode: адамс: нужно не менее 4 шагов, есть %d", n)
	}
	us := make([][]float64, n+1)
	fs := make([][]float64, n+1)
	u := clone(u0)
	us[0] = u

	for k := 0; k < 3; k++ {
		u, err = rk4Step(f, xs[k], u, h)
		if err != nil {
			return nil, nil, err
		}
		us[k+1] = clone(u)
	}
	for k := 0; k <= 3; k++ {
		fs[k], err = derivative(f, xs[k], us[k])
		if err != nil {
			return nil, nil, err
		}
	}

	for k := 3; k < n; k++ {
		du := axpy(axpy(axpy(scale(fs[k], 55), fs[k-1], -59), fs[k-2], 37), fs[k-3], -9)
		u = axpy(u, du, h/24)
		if !finite(u) {
			return nil, nil, fmt.Errorf("ode: адамс: решение расходится при x = %g", xs[k+1])
		}
		us[k+1] = clone(u)
		fs[k+1], err = derivative(f, xs[k+1], u)
		if err != nil {
			return nil, nil, err
		}
	}
	return xs, us, nil
}

func nodes(x0, xk, h float64) ([]float64, int, error) {
	if !finite([]float64{x0, xk, h}) || h <= 0 || xk <= x0 {
		return nil, 0, fmt.Errorf("ode: некорректные параметры: x0 = %g, xk = %g, h = %g", x0, xk, h)
	}
	steps := (xk - x0) / h
	rounded := math.Round(steps)
	if !finite([]float64{steps}) || rounded >= float64(int(^uint(0)>>1)) || math.Abs(steps-rounded) > 1e-9 {
		return nil, 0, fmt.Errorf("ode: шаг %g не делит отрезок [%g, %g] или сетка слишком велика", h, x0, xk)
	}
	n := int(rounded)
	if n < 1 {
		return nil, 0, fmt.Errorf("ode: шаг %g слишком велик для отрезка [%g, %g]", h, x0, xk)
	}
	xs := make([]float64, n+1)
	for i := range xs {
		xs[i] = x0 + float64(i)*h
		if i == n {
			xs[i] = xk
		}
		if !finite(xs[i:i+1]) || (i > 0 && xs[i] <= xs[i-1]) {
			return nil, 0, fmt.Errorf("ode: непредставимый узел сетки %d", i)
		}
	}
	return xs, n, nil
}

func derivative(f Func, x float64, u []float64) ([]float64, error) {
	if !finite([]float64{x}) || !finite(u) {
		return nil, fmt.Errorf("ode: неконечное состояние при x = %g", x)
	}
	du := f(x, u)
	if len(du) != len(u) || !finite(du) {
		return nil, fmt.Errorf("ode: неверная размерность или неконечная производная при x = %g: нужно %d, получено %d", x, len(u), len(du))
	}
	return clone(du), nil
}

func clone(u []float64) []float64 {
	v := make([]float64, len(u))
	copy(v, u)
	return v
}

func axpy(a, b []float64, c float64) []float64 {
	v := make([]float64, len(a))
	for i := range a {
		v[i] = a[i] + c*b[i]
	}
	return v
}

func scale(a []float64, c float64) []float64 {
	v := make([]float64, len(a))
	for i := range a {
		v[i] = c * a[i]
	}
	return v
}

func finite(u []float64) bool {
	for _, v := range u {
		if math.IsNaN(v) || math.IsInf(v, 0) {
			return false
		}
	}
	return true
}
