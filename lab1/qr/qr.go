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

// Pair — комплексно-сопряжённая пара собственных значений.
type Pair struct{ Re, Im float64 }

func Find(a *matrix.Matrix, maxIter int) ([]float64, []Pair, error) {
	n := a.Rows()
	if a.Cols() != n {
		return nil, nil, fmt.Errorf("qr: матрица не квадратная (%dx%d)", n, a.Cols())
	}
	b := a.Copy()
	scale := 1 + b.Norm2()

	for it := 0; it < maxIter; it++ {
		q, r, err := Factor(b)
		if err != nil {
			return nil, nil, err
		}
		b, err = r.Multiply(q)
		if err != nil {
			return nil, nil, err
		}
		if maxSubDiag(b) <= matrix.Eps*scale {
			break
		}
	}

	tol := matrix.Eps * scale
	var reals []float64
	var pairs []Pair
	for i := 0; i < n; {
		sub := 0.0
		if i+1 < n {
			sub = math.Abs(b.Get(i+1, i))
		}
		if sub <= tol {
			reals = append(reals, b.Get(i, i))
			i++
			continue
		}
		// блок 2×2 с комплексными СЗ: λ = m ± √q, q = ((a−d)/2)² + bc
		a11, a12 := b.Get(i, i), b.Get(i, i+1)
		a21, a22 := b.Get(i+1, i), b.Get(i+1, i+1)
		m := (a11 + a22) / 2
		q := (a11-a22)*(a11-a22)/4 + a12*a21
		if q > 0 {
			// вещественные СЗ блока (блок ещё не сошёлся)
			s := math.Sqrt(q)
			reals = append(reals, m+s, m-s)
		} else {
			im := math.Sqrt(-q)
			pairs = append(pairs, Pair{m, im}, Pair{m, -im})
		}
		i += 2
	}
	return reals, pairs, nil
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
