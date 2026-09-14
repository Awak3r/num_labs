package iterative

import (
	"fmt"
	"math"

	"lab1/matrix"
)

func SimpleIterations(a, b *matrix.Matrix, eps float64, maxIter int) (*matrix.Matrix, int, error) {
	n := a.Rows()
	if err := checkDiagonal(a); err != nil {
		return nil, 0, err
	}
	x := matrix.Zeros(n, 1)
	for k := 0; k < maxIter; k++ {
		xNew := matrix.Zeros(n, 1)
		for i := 0; i < n; i++ {
			var s float64
			for j := 0; j < n; j++ {
				if j != i {
					s += a.Get(i, j) * x.Get(j, 0)
				}
			}
			xNew.Set(i, 0, (b.Get(i, 0)-s)/a.Get(i, i))
		}
		diff, _ := xNew.Sub(x)
		if diff.MaxAbs() < eps {
			return xNew, k + 1, nil
		}
		x = xNew
	}
	return nil, 0, fmt.Errorf("iterative: простые итерации не сошлись за %d итераций", maxIter)
}

func Seidel(a, b *matrix.Matrix, eps float64, maxIter int) (*matrix.Matrix, int, error) {
	n := a.Rows()
	if err := checkDiagonal(a); err != nil {
		return nil, 0, err
	}
	x := matrix.Zeros(n, 1)
	for k := 0; k < maxIter; k++ {
		xNew := x.Copy()
		for i := 0; i < n; i++ {
			var s float64
			for j := 0; j < i; j++ {
				s += a.Get(i, j) * xNew.Get(j, 0)
			}
			for j := i + 1; j < n; j++ {
				s += a.Get(i, j) * x.Get(j, 0)
			}
			xNew.Set(i, 0, (b.Get(i, 0)-s)/a.Get(i, i))
		}
		diff, _ := xNew.Sub(x)
		if diff.MaxAbs() < eps {
			return xNew, k + 1, nil
		}
		x = xNew
	}
	return nil, 0, fmt.Errorf("iterative: метод Зейделя не сошёлся за %d итераций", maxIter)
}

func checkDiagonal(a *matrix.Matrix) error {
	for i := 0; i < a.Rows(); i++ {
		if math.Abs(a.Get(i, i)) < matrix.Eps {
			return fmt.Errorf("iterative: нулевой диагональный элемент A[%d][%d]", i, i)
		}
	}
	return nil
}
