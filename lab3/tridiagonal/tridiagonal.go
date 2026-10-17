package tridiagonal

import (
	"fmt"
	"math"

	"lab3/matrix"
)

// Solve решает трёхдиагональную систему методом прогонки (Томас).
// main — главная диагональ (n), super — наддиагональ (n-1), sub — поддиагональ (n-1),
// d — правая часть (n). Возвращает решение x.
func Solve(main, super, sub, d []float64) ([]float64, error) {
	n := len(main)
	if len(super) != n-1 || len(sub) != n-1 || len(d) != n {
		return nil, fmt.Errorf("tridiagonal: несовместимые размеры: main=%d, super=%d, sub=%d, d=%d", n, len(super), len(sub), len(d))
	}
	if n == 0 {
		return nil, fmt.Errorf("tridiagonal: пустая система")
	}
	// Прямой ход: α, β
	alpha := make([]float64, n)
	beta := make([]float64, n)
	if math.Abs(main[0]) < matrix.Eps {
		return nil, fmt.Errorf("tridiagonal: нулевой опорный элемент на шаге 0")
	}
	alpha[0] = -super[0] / main[0]
	beta[0] = d[0] / main[0]
	for i := 1; i < n; i++ {
		denom := main[i] + sub[i-1]*alpha[i-1]
		if math.Abs(denom) < matrix.Eps {
			return nil, fmt.Errorf("tridiagonal: нулевой опорный элемент на шаге %d", i)
		}
		if i < n-1 {
			alpha[i] = -super[i] / denom
		} else {
			alpha[i] = 0
		}
		beta[i] = (d[i] - sub[i-1]*beta[i-1]) / denom
	}
	// Обратный ход
	x := make([]float64, n)
	x[n-1] = beta[n-1]
	for i := n - 2; i >= 0; i-- {
		x[i] = alpha[i]*x[i+1] + beta[i]
	}
	return x, nil
}


