package funcs

import (
	"fmt"
	"math"

	"lab2/matrix"
)

func F(x float64) float64 {
	return math.Sqrt(1-x*x) - math.Exp(x) + 0.1
}

func DF(x float64) float64 {
	return -x/math.Sqrt(1-x*x) - math.Exp(x)
}

func Phi(x float64) float64 {
	return math.Log(math.Sqrt(1-x*x) + 0.1)
}

func DPhi(x float64) float64 {
	return -x / (math.Sqrt(1-x*x) * (math.Sqrt(1-x*x) + 0.1))
}

func A() float64 { return 4 }

func F1(x *matrix.Matrix) float64 {
	a := A()
	return (x.Get(0, 0)*x.Get(0, 0) + a*a) * x.Get(1, 0) - a*a*a
}

func F2(x *matrix.Matrix) float64 {
	a := A()
	x1, x2 := x.Get(0, 0), x.Get(1, 0)
	return (x1-a/2)*(x1-a/2) + (x2-a/2)*(x2-a/2) - a*a
}

func PhiSystem(x *matrix.Matrix) *matrix.Matrix {
	a := A()
	t := a*a - (x.Get(1, 0)-a/2)*(x.Get(1, 0)-a/2)
	if t < 0 {
		return nil
	}
	x1 := a/2 + math.Sqrt(t)
	x2 := a * a * a / (x1*x1 + a*a)
	if math.IsNaN(x2) || math.IsInf(x2, 0) {
		return nil
	}
	return matrix.Column([]float64{x1, x2})
}

func FSystem(x *matrix.Matrix) *matrix.Matrix {
	if x.Get(1, 0) <= 0 {
		return nil
	}
	return matrix.Column([]float64{F1(x), F2(x)})
}

func Jacobian(x *matrix.Matrix) *matrix.Matrix {
	a := A()
	x1, x2 := x.Get(0, 0), x.Get(1, 0)
	return matrix.Create([][]float64{
		{2 * x1 * x2, x1*x1 + a*a},
		{2 * (x1 - a/2), 2 * (x2 - a/2)},
	})
}

func VectorString(x *matrix.Matrix) string {
	return fmt.Sprintf("(%.*g, %.*g)", 6, x.Get(0, 0), 6, x.Get(1, 0))
}
