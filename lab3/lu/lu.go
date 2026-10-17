package lu

import (
	"fmt"
	"math"

	"lab3/matrix"
)

func Decompose(a *matrix.Matrix) (*matrix.Matrix, *matrix.Matrix, *matrix.Matrix, error) {
	if a.Rows() != a.Cols() {
		return nil, nil, nil, fmt.Errorf("lu: матрица должна быть квадратной, получено %dx%d", a.Rows(), a.Cols())
	}
	n := a.Rows()
	u := a.Copy()
	l := matrix.MakeIdentityMatrix(n)
	p := matrix.MakeIdentityMatrix(n)

	for k := 0; k < n; k++ {
		pivot := k
		max := math.Abs(u.Get(k, k))
		for i := k + 1; i < n; i++ {
			if v := math.Abs(u.Get(i, k)); v > max {
				max = v
				pivot = i
			}
		}
		if max < matrix.Eps {
			return nil, nil, nil, fmt.Errorf("lu: матрица вырождена (нулевой столбец %d)", k)
		}
		if pivot != k {
			u.SwapRows(pivot, k)
			p.SwapRows(pivot, k)
			for j := 0; j < k; j++ {
				tmp := l.Get(pivot, j)
				l.Set(pivot, j, l.Get(k, j))
				l.Set(k, j, tmp)
			}
		}
		for i := k + 1; i < n; i++ {
			l.Set(i, k, u.Get(i, k)/u.Get(k, k))
			for j := k; j < n; j++ {
				u.Set(i, j, u.Get(i, j)-l.Get(i, k)*u.Get(k, j))
			}
		}
	}
	return l, u, p, nil
}

func Solve(a, b *matrix.Matrix) (*matrix.Matrix, error) {
	if a.Rows() != a.Cols() {
		return nil, fmt.Errorf("lu: матрица должна быть квадратной")
	}
	if b.Rows() != a.Rows() || b.Cols() != 1 {
		return nil, fmt.Errorf("lu: b должен быть вектором-столбцом %dx1", a.Rows())
	}
	l, u, p, err := Decompose(a)
	if err != nil {
		return nil, err
	}
	pb, err := p.Multiply(b)
	if err != nil {
		return nil, err
	}
	return backward(u, forward(l, pb))
}

func forward(l, pb *matrix.Matrix) *matrix.Matrix {
	n := pb.Rows()
	y := matrix.Zeros(n, 1)
	for i := 0; i < n; i++ {
		sum := pb.Get(i, 0)
		for j := 0; j < i; j++ {
			sum -= l.Get(i, j) * y.Get(j, 0)
		}
		y.Set(i, 0, sum)
	}
	return y
}

func backward(u, y *matrix.Matrix) (*matrix.Matrix, error) {
	n := y.Rows()
	x := matrix.Zeros(n, 1)
	for i := n - 1; i >= 0; i-- {
		sum := y.Get(i, 0)
		for j := i + 1; j < n; j++ {
			sum -= u.Get(i, j) * x.Get(j, 0)
		}
		diag := u.Get(i, i)
		if math.Abs(diag) < matrix.Eps {
			return nil, fmt.Errorf("lu: нулевой диагональный элемент U[%d][%d]", i, i)
		}
		x.Set(i, 0, sum/diag)
	}
	return x, nil
}
