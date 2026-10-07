package jacobi

import (
	"fmt"
	"math"

	"lab1/matrix"
)

func Solve(a *matrix.Matrix, eps float64, maxIter int) ([]float64, *matrix.Matrix, int, error) {
	n := a.Rows()
	if a.Cols() != n {
		return nil, nil, 0, fmt.Errorf("jacobi: матрица не квадратная (%dx%d)", n, a.Cols())
	}
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if math.Abs(a.Get(i, j)-a.Get(j, i)) > matrix.Eps {
				return nil, nil, 0, fmt.Errorf("jacobi: матрица не симметрична: a[%d][%d]=%g, a[%d][%d]=%g", i, j, a.Get(i, j), j, i, a.Get(j, i))
			}
		}
	}
	work := a.Copy()
	v := matrix.MakeIdentityMatrix(n)
	var iters int
	for iters = 0; iters < maxIter; iters++ {
		if offDiagNorm(work) < eps {
			break
		}
		p, q, _ := maxOffDiag(work)
		theta := (work.Get(q, q) - work.Get(p, p)) / (2 * work.Get(p, q))
		t := 1.0
		if theta != 0 {
			t = math.Copysign(1, theta) / (math.Abs(theta) + math.Sqrt(theta*theta+1))
		}
		c := 1 / math.Sqrt(t*t+1)
		s := c * t
		rotateA(work, p, q, c, s)
		rotateV(v, p, q, c, s)
	}
	if iters == maxIter {
		return nil, nil, 0, fmt.Errorf("jacobi: не сошёлся за %d вращений", maxIter)
	}
	vals := make([]float64, n)
	for i := 0; i < n; i++ {
		vals[i] = work.Get(i, i)
	}
	return vals, v, iters, nil
}

func maxOffDiag(a *matrix.Matrix) (int, int, float64) {
	p, q := 0, 1
	var m float64
	for i := 0; i < a.Rows(); i++ {
		for j := i + 1; j < a.Cols(); j++ {
			if v := math.Abs(a.Get(i, j)); v > m {
				m, p, q = v, i, j
			}
		}
	}
	return p, q, m
}

func offDiagNorm(a *matrix.Matrix) float64 {
	var s float64
	for i := 1; i < a.Rows(); i++ {
		for j := 0; j < i; j++ {
			s += a.Get(i, j) * a.Get(i, j)
		}
	}
	return math.Sqrt(s)
}

func rotateA(a *matrix.Matrix, p, q int, c, s float64) {
	app := a.Get(p, p)
	aqq := a.Get(q, q)
	apq := a.Get(p, q)
	for k := 0; k < a.Rows(); k++ {
		if k == p || k == q {
			continue
		}
		apk := a.Get(p, k)
		aqk := a.Get(q, k)
		a.Set(p, k, c*apk-s*aqk)
		a.Set(q, k, s*apk+c*aqk)
		a.Set(k, p, a.Get(p, k))
		a.Set(k, q, a.Get(q, k))
	}
	a.Set(p, p, c*c*app-2*s*c*apq+s*s*aqq)
	a.Set(q, q, s*s*app+2*s*c*apq+c*c*aqq)
	a.Set(p, q, 0)
	a.Set(q, p, 0)
}

func rotateV(v *matrix.Matrix, p, q int, c, s float64) {
	for r := 0; r < v.Rows(); r++ {
		vp := v.Get(r, p)
		vq := v.Get(r, q)
		v.Set(r, p, c*vp-s*vq)
		v.Set(r, q, s*vp+c*vq)
	}
}
