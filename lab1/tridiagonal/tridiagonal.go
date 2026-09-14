package tridiagonal

import (
	"fmt"
	"math"

	"lab1/matrix"
)

// Solve решает трёхдиагональную систему a*x[i-1] + b*x[i] + c*x[i+1] = d
// методом прогонки (Thomas). Хранит только диагонали.
func Solve(main, super, sub, d []float64) (*matrix.Matrix, []float64, []float64, error) {
	n := len(main)
	if n == 0 {
		return nil, nil, nil, fmt.Errorf("tridiagonal: пустая система")
	}
	if n == 1 {
		if math.Abs(main[0]) < matrix.Eps {
			return nil, nil, nil, fmt.Errorf("tridiagonal: нулевой коэффициент main[0]")
		}
		x := []float64{d[0] / main[0]}
		return matrix.Column(x), nil, nil, nil
	}

	alpha := make([]float64, n)
	beta := make([]float64, n)

	if math.Abs(main[0]) < matrix.Eps {
		return nil, nil, nil, fmt.Errorf("tridiagonal: нулевой коэффициент main[0]")
	}
	alpha[0] = -super[0] / main[0]
	beta[0] = d[0] / main[0]

	for i := 1; i < n; i++ {
		denom := main[i] + sub[i-1]*alpha[i-1]
		if math.Abs(denom) < matrix.Eps {
			return nil, nil, nil, fmt.Errorf("tridiagonal: знаменатель близок к нулю на шаге %d", i)
		}
		if i < n-1 {
			alpha[i] = -super[i] / denom
		}
		beta[i] = (d[i] - sub[i-1]*beta[i-1]) / denom
	}

	x := make([]float64, n)
	x[n-1] = beta[n-1]
	for i := n - 2; i >= 0; i-- {
		x[i] = alpha[i]*x[i+1] + beta[i]
	}
	return matrix.Column(x), alpha, beta, nil
}
