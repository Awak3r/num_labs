package lsm

import (
	"fmt"
	"strings"

	"lab3/lu"
	"lab3/matrix"
)

// Fit строит аппроксимирующий многочлен степени degree МНК
// через нормальную систему A^T·A·c = A^T·y (базис 1, x, ..., x^degree).
func Fit(x, y []float64, degree int) (Poly, error) {
	if degree < 1 {
		return Poly{}, fmt.Errorf("lsm: степень должна быть ≥1, есть %d", degree)
	}
	if len(x) != len(y) || len(x) < degree+1 {
		return Poly{}, fmt.Errorf("lsm: точек %d недостаточно для степени %d", len(x), degree)
	}
	m := degree + 1
	// Матрица A: A[i][j] = x_i^j
	normal := matrix.Zeros(m, m)
	rhs := matrix.Column(make([]float64, m))
	for i := 0; i < m; i++ {
		for j := 0; j < m; j++ {
			s := 0.0
			for k := range x {
				s += power(x[k], i+j)
			}
			normal.Set(i, j, s)
		}
		s := 0.0
		for k := range x {
			s += y[k] * power(x[k], i)
		}
		rhs.Set(i, 0, s)
	}
	coef, err := lu.Solve(normal, rhs)
	if err != nil {
		return Poly{}, fmt.Errorf("lsm: %w", err)
	}
	c := make([]float64, m)
	for i := 0; i < m; i++ {
		c[i] = coef.Get(i, 0)
	}
	return Poly{Coef: c}, nil
}

type Poly struct {
	Coef []float64 // по возрастанию степени
}

func (p Poly) Eval(x float64) float64 {
	res := 0.0
	for i := len(p.Coef) - 1; i >= 0; i-- {
		res = res*x + p.Coef[i]
	}
	return res
}

// SSE — сумма квадратов ошибок по всем точкам.
func (p Poly) SSE(x, y []float64) float64 {
	s := 0.0
	for i := range x {
		d := p.Eval(x[i]) - y[i]
		s += d * d
	}
	return s
}

// Form — аналитический вид строки.
func (p Poly) Form() string {
	var parts []string
	for i := len(p.Coef) - 1; i >= 0; i-- {
		c := p.Coef[i]
		if i == 0 {
			parts = append(parts, fmt.Sprintf("%+.*g", 6, c))
		} else if i == 1 {
			parts = append(parts, fmt.Sprintf("%+.*g·x", 6, c))
		} else {
			parts = append(parts, fmt.Sprintf("%+.*g·x^%d", 6, c, i))
		}
	}
	return strings.Join(parts, " ")
}

func power(base float64, e int) float64 {
	r := 1.0
	for i := 0; i < e; i++ {
		r *= base
	}
	return r
}
