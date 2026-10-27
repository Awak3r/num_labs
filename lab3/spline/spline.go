package spline

import (
	"fmt"

	"lab3/tridiagonal"
)

type Piece struct {
	A, B, C, D float64
	X0, X1     float64
	Index      int
}

type Spline struct {
	X, Y    []float64
	Pieces  []Piece
	H       []float64
}

func Build(x, y []float64) (*Spline, error) {
	n := len(x)
	if n < 2 || len(y) != n {
		return nil, fmt.Errorf("spline: нужно ≥2 узлов и |y|=|x|, есть n=%d, len(y)=%d", n, len(y))
	}
	for i := 1; i < n; i++ {
		if x[i]-x[i-1] <= 0 {
			return nil, fmt.Errorf("spline: узлы должны быть строго возрастающими")
		}
	}
	h := make([]float64, n-1)
	for i := 0; i < n-1; i++ {
		h[i] = x[i+1] - x[i]
	}

	m := n - 2
	main := make([]float64, m)
	super := make([]float64, m-1)
	sub := make([]float64, m-1)
	d := make([]float64, m)
	for i := 1; i <= m; i++ {
		main[i-1] = 2 * (h[i-1] + h[i])
		d[i-1] = 6 * ((y[i+1]-y[i])/h[i] - (y[i]-y[i-1])/h[i-1])
		if i > 1 {
			sub[i-2] = h[i-1]
		}
		if i < m {
			super[i-1] = h[i]
		}
	}
	cSol, err := tridiagonal.Solve(main, super, sub, d)
	if err != nil {
		return nil, fmt.Errorf("spline: %w", err)
	}
	c := make([]float64, n)
	copy(c[1:m+1], cSol)

	b := make([]float64, n)
	a := make([]float64, n)
	dd := make([]float64, n)
	for i := 1; i < n; i++ {
		a[i] = y[i-1]
		b[i] = (y[i]-y[i-1])/h[i-1] - h[i-1]*(c[i]+2*c[i-1])/3
		dd[i] = (c[i] - c[i-1]) / (3 * h[i-1])
	}

	pieces := make([]Piece, n-1)
	for i := 1; i < n; i++ {
		pieces[i-1] = Piece{
			A: a[i], B: b[i], C: c[i-1], D: dd[i],
			X0: x[i-1], X1: x[i], Index: i,
		}
	}
	return &Spline{X: x, Y: y, Pieces: pieces, H: h}, nil
}

func (s *Spline) Evaluate(x float64) (float64, int, error) {
	if x < s.X[0]-1e-12 || x > s.X[len(s.X)-1]+1e-12 {
		return 0, 0, fmt.Errorf("spline: x=%g вне [%g, %g]", x, s.X[0], s.X[len(s.X)-1])
	}
	for i := 1; i < len(s.X); i++ {
		if x <= s.X[i]+1e-12 {
			p := s.Pieces[i-1]
			t := x - p.X0
			return p.A + p.B*t + p.C*t*t + p.D*t*t*t, p.Index, nil
		}
	}
	p := s.Pieces[len(s.Pieces)-1]
	t := x - p.X0
	return p.A + p.B*t + p.C*t*t + p.D*t*t*t, p.Index, nil
}

func (s *Spline) SecondDerivative(x float64) (float64, error) {
	if x < s.X[0]-1e-12 || x > s.X[len(s.X)-1]+1e-12 {
		return 0, fmt.Errorf("spline: x=%g вне диапазона", x)
	}
	for i := 1; i < len(s.X); i++ {
		if x <= s.X[i]+1e-12 {
			p := s.Pieces[i-1]
			t := x - p.X0
			return 2*p.C + 6*p.D*t, nil
		}
	}
	p := s.Pieces[len(s.Pieces)-1]
	t := x - p.X0
	return 2*p.C + 6*p.D*t, nil
}
