package qr

import (
	"fmt"
	"math"

	"lab1/matrix"
)

func Factor(a *matrix.Matrix) (*matrix.Matrix, *matrix.Matrix, error) {
	n := a.Rows()
	if a.Cols() != n {
		return nil, nil, fmt.Errorf("qr: матрица не квадратная (%dx%d)", n, a.Cols())
	}
	r := a.Copy()
	q := matrix.MakeIdentityMatrix(n)

	for k := 0; k < n-1; k++ {
		var norm float64
		for i := k; i < n; i++ {
			norm += r.Get(i, k) * r.Get(i, k)
		}
		norm = math.Sqrt(norm)
		if norm <= matrix.Eps {
			continue
		}
		alfa := -math.Copysign(norm, r.Get(k, k))

		v := make([]float64, n)
		for i := k; i < n; i++ {
			v[i] = r.Get(i, k)
		}
		v[k] -= alfa

		var vv float64
		for i := k; i < n; i++ {
			vv += v[i] * v[i]
		}
		tau := 2.0 / vv

		for j := k; j < n; j++ {
			var dot float64
			for i := k; i < n; i++ {
				dot += v[i] * r.Get(i, j)
			}
			for i := k; i < n; i++ {
				r.Set(i, j, r.Get(i, j)-tau*v[i]*dot)
			}
		}
		r.Set(k, k, alfa)
		for i := k + 1; i < n; i++ {
			r.Set(i, k, 0)
		}

		for row := 0; row < n; row++ {
			var dot float64
			for i := k; i < n; i++ {
				dot += v[i] * q.Get(row, i)
			}
			for i := k; i < n; i++ {
				q.Set(row, i, q.Get(row, i)-tau*v[i]*dot)
			}
		}
	}
	return q, r, nil
}

func Find(a *matrix.Matrix, maxIter int) ([]float64, bool, error) {
	n := a.Rows()
	if a.Cols() != n {
		return nil, false, fmt.Errorf("qr: матрица не квадратная (%dx%d)", n, a.Cols())
	}
	b := a.Copy()
	scale := 1 + b.Norm2()

	for it := 0; it < maxIter; it++ {
		q, r, err := Factor(b)
		if err != nil {
			return nil, false, err
		}
		b, err = r.Multiply(q)
		if err != nil {
			return nil, false, err
		}
		if maxSubDiag(b) <= matrix.Eps*scale {
			break
		}
	}

	tol := matrix.Eps * scale
	hasPair := maxSubDiag(b) > tol
	var real []float64
	for i := 0; i < n; {
		sub := 0.0
		if i+1 < n {
			sub = math.Abs(b.Get(i+1, i))
		}
		if sub <= tol {
			real = append(real, b.Get(i, i))
			i++
		} else {
			i += 2
		}
	}
	return real, hasPair, nil
}

func maxSubDiag(a *matrix.Matrix) float64 {
	var m float64
	for i := 1; i < a.Rows(); i++ {
		if v := math.Abs(a.Get(i, i-1)); v > m {
			m = v
		}
	}
	return m
}
