package funcs

import "math"

// ---- 4.1 Задача Коши (P_4_1.doc, вариант 3) ----
// y'' - 2y - 4x^2·e^(x^2) = 0  <=>  y'' = 2y + 4x^2·e^(x^2),
// y(0) = 3, y'(0) = 0, x ∈ [0,1], h = 0.1.
// Точное решение: y(x) = e^(x^2) + e^(x√2) + e^(−x√2).

// FCauchy — правая часть системы u' = f(x, u), u = (y, y').
func FCauchy(x float64, u []float64) []float64 {
	return []float64{u[1], 2*u[0] + 4*x*x*math.Exp(x*x)}
}

// CauchyU0 — начальные условия y(0) = 3, y'(0) = 0.
func CauchyU0() []float64 { return []float64{3, 0} }

// ExactCauchy — точное решение задачи Коши.
func ExactCauchy(x float64) float64 {
	return math.Exp(x*x) + math.Exp(x*math.Sqrt2) + math.Exp(-x*math.Sqrt2)
}

// ---- 4.2 Краевая задача (P_4_2.doc, вариант 3) ----
// x^2·(x+1)·y'' - 2y = 0  <=>  y'' + p(x)·y' + q(x)·y = r(x),
// p(x) = 0, q(x) = −2/(x^2·(x+1)), r(x) = 0.
// Граничные условия: y'(1) = −1, 2y(2) − 4y'(2) = 4, x ∈ [1,2].
// Точное решение: y(x) = 1 + 1/x.

// FBoundary — правая часть системы u' = f(x, u), u = (y, y') для стрельбы.
func FBoundary(x float64, u []float64) []float64 {
	return []float64{u[1], 2 * u[0] / (x * x * (x + 1))}
}

// PBoundary/QBoundary/RBoundary — коэффициенты линейного уравнения.
func PBoundary(float64) float64 { return 0 }
func QBoundary(x float64) float64 {
	return -2 / (x * x * (x + 1))
}
func RBoundary(float64) float64 { return 0 }

// BC0 — левое краевое условие a0·y(x0) + b0·y'(x0) = g0 (y'(1) = −1).
func BC0() (a0, b0, g0 float64) { return 0, 1, -1 }

// BC1 — правое краевое условие a1·y(xk) + b1·y'(xk) = g1 (2y(2) − 4y'(2) = 4).
func BC1() (a1, b1, g1 float64) { return 2, -4, 4 }

// DU0 — известная производная в левой точке для стрельбы: y'(x0) = −1.
func DU0() float64 { return -1 }

// RightResidual — невязка правого краевого условия: 2y − 4y' − 4.
func RightResidual(u []float64) float64 { return 2*u[0] - 4*u[1] - 4 }

// ShootS0/ShootS1 — начальные приближения s = y(1) (графически y(1) ≈ 2).
func ShootS0() float64 { return 2 }
func ShootS1() float64 { return 2.5 }

// ExactBoundary — точное решение краевой задачи.
func ExactBoundary(x float64) float64 { return 1 + 1/x }
