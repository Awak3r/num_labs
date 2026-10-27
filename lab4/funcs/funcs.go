package funcs

import "math"

func FCauchy(x float64, u []float64) []float64 {
	return []float64{u[1], 2*u[0] + 4*x*x*math.Exp(x*x)}
}

func CauchyU0() []float64 { return []float64{3, 0} }

func ExactCauchy(x float64) float64 {
	return math.Exp(x*x) + math.Exp(x*math.Sqrt2) + math.Exp(-x*math.Sqrt2)
}

func FBoundary(x float64, u []float64) []float64 {
	return []float64{u[1], 2 * u[0] / (x * x * (x + 1))}
}

func PBoundary(float64) float64 { return 0 }
func QBoundary(x float64) float64 {
	return -2 / (x * x * (x + 1))
}
func RBoundary(float64) float64 { return 0 }

func BC0() (a0, b0, g0 float64) { return 0, 1, -1 }

func BC1() (a1, b1, g1 float64) { return 2, -4, 4 }

func DU0() float64 { return -1 }

func RightResidual(u []float64) float64 { return 2*u[0] - 4*u[1] - 4 }

func ShootS0() float64 { return 2 }
func ShootS1() float64 { return 2.5 }

func ExactBoundary(x float64) float64 { return 1 + 1/x }
