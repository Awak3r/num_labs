package ode

import (
	"fmt"
	"math"
)

// Func — правая часть системы первого порядка u' = f(x, u).
type Func func(x float64, u []float64) []float64

// Euler — явный метод Эйлера (p = 1).
func Euler(f Func, x0, xk, h float64, u0 []float64) ([]float64, [][]float64, error) {
	xs, n, err := nodes(x0, xk, h)
	if err != nil {
		return nil, nil, err
	}
	us := make([][]float64, n+1)
	u := clone(u0)
	us[0] = u
	for k := 0; k < n; k++ {
		du := f(xs[k], u)
		if du == nil {
			return nil, nil, fmt.Errorf("ode: эйлер: f вне области определения при x = %g", xs[k])
		}
		u = axpy(u, du, h)
		if !finite(u) {
			return nil, nil, fmt.Errorf("ode: эйлер: решение расходится при x = %g", xs[k+1])
		}
		us[k+1] = clone(u)
	}
	return xs, us, nil
}

// RK4 — классический метод Рунге-Кутты 4-го порядка (p = 4).
func RK4(f Func, x0, xk, h float64, u0 []float64) ([]float64, [][]float64, error) {
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
	k1 := f(x, u)
	if k1 == nil {
		return nil, fmt.Errorf("ode: рк4: f вне области определения при x = %g", x)
	}
	k2 := f(x+h/2, axpy(u, k1, h/2))
	if k2 == nil {
		return nil, fmt.Errorf("ode: рк4: f вне области определения при x = %g", x+h/2)
	}
	k3 := f(x+h/2, axpy(u, k2, h/2))
	if k3 == nil {
		return nil, fmt.Errorf("ode: рк4: f вне области определения при x = %g", x+h/2)
	}
	k4 := f(x+h, axpy(u, k3, h))
	if k4 == nil {
		return nil, fmt.Errorf("ode: рк4: f вне области определения при x = %g", x+h)
	}
	sum := axpy(axpy(axpy(k1, k2, 2), k3, 2), k4, 1)
	u = axpy(u, sum, h/6)
	if !finite(u) {
		return nil, fmt.Errorf("ode: рк4: решение расходится при x = %g", x+h)
	}
	return u, nil
}

// Adams4 — явный метод Адамса 4-го порядка (p = 4); разгон — РК4.
func Adams4(f Func, x0, xk, h float64, u0 []float64) ([]float64, [][]float64, error) {
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
	// разгон: первые 3 шага РК4
	for k := 0; k < 3; k++ {
		u, err = rk4Step(f, xs[k], u, h)
		if err != nil {
			return nil, nil, err
		}
		us[k+1] = clone(u)
	}
	for k := 0; k <= 3; k++ {
		fs[k] = f(xs[k], us[k])
		if fs[k] == nil {
			return nil, nil, fmt.Errorf("ode: адамс: f вне области определения при x = %g", xs[k])
		}
	}
	// шаги Адамса
	for k := 3; k < n; k++ {
		du := axpy(axpy(axpy(scale(fs[k], 55), fs[k-1], -59), fs[k-2], 37), fs[k-3], -9)
		u = axpy(u, du, h/24)
		if !finite(u) {
			return nil, nil, fmt.Errorf("ode: адамс: решение расходится при x = %g", xs[k+1])
		}
		us[k+1] = clone(u)
		fs[k+1] = f(xs[k+1], u)
		if fs[k+1] == nil {
			return nil, nil, fmt.Errorf("ode: адамс: f вне области определения при x = %g", xs[k+1])
		}
	}
	return xs, us, nil
}

// nodes — узлы сетки x_i = x0 + i·h, n — число шагов.
func nodes(x0, xk, h float64) ([]float64, int, error) {
	if h <= 0 || xk <= x0 {
		return nil, 0, fmt.Errorf("ode: некорректные параметры: x0 = %g, xk = %g, h = %g", x0, xk, h)
	}
	n := int(math.Round((xk - x0) / h))
	if n < 1 {
		return nil, 0, fmt.Errorf("ode: шаг %g слишком велик для отрезка [%g, %g]", h, x0, xk)
	}
	xs := make([]float64, n+1)
	for i := range xs {
		xs[i] = x0 + float64(i)*h
	}
	return xs, n, nil
}

func clone(u []float64) []float64 {
	v := make([]float64, len(u))
	copy(v, u)
	return v
}

// axpy — a + c·b (новый срез).
func axpy(a, b []float64, c float64) []float64 {
	v := make([]float64, len(a))
	for i := range a {
		v[i] = a[i] + c*b[i]
	}
	return v
}

// scale — c·a (новый срез).
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
