package demo

import (
	"fmt"

	"lab3/input"
	"lab3/spline"
)

func Spline() {
	fmt.Println("=== 3.2 Кубический сплайн, x* = 1.5 ===")
	pts, err := input.ReadPointsXstar("data/spline.txt")
	if err != nil {
		fmt.Println("ошибка:", err)
		return
	}
	fmt.Println("исходные узлы:")
	for i := range pts.X {
		fmt.Printf("  x%d = %8.4f   y%d = %12.6f\n", i, pts.X[i], i, pts.Y[i])
	}
	sp, err := spline.Build(pts.X, pts.Y)
	if err != nil {
		fmt.Println("ошибка:", err)
		return
	}
	fmt.Println("коэффициенты сплайна на каждом участке S_i(x) = a_i + b_i(x−x_{i−1}) + c_i(x−x_{i−1})² + d_i(x−x_{i−1})³:")
	for _, p := range sp.Pieces {
		fmt.Printf("  [%g, %g]: a=%12.6f b=%12.6f c=%12.6f d=%12.6f\n", p.X0, p.X1, p.A, p.B, p.C, p.D)
	}
	v, idx, err := sp.Evaluate(pts.XStar)
	if err != nil {
		fmt.Println("ошибка:", err)
		return
	}
	fmt.Printf("S(x*) = S_%d(%g) = %12.6f\n", idx, pts.XStar, v)
	fmt.Println("проверка прохождения через узлы:")
	ok := true
	for i := range sp.X {
		si, _, err := sp.Evaluate(sp.X[i])
		if err != nil {
			fmt.Println("ошибка:", err)
			return
		}
		diff := si - sp.Y[i]
		if diff > 1e-9 || diff < -1e-9 {
			ok = false
		}
		fmt.Printf("  S(x%d) = %12.6f  y%d = %12.6f  |Δ| = %g\n", i, si, i, sp.Y[i], abs(diff))
	}
	if ok {
		fmt.Println("сплайн проходит через все узлы: да")
	} else {
		fmt.Println("сплайн проходит через все узлы: НЕТ")
	}
	fmt.Println("проверка гладкости во внутренних узлах:")
	for i := 1; i < len(sp.Pieces); i++ {
		left, right := sp.Pieces[i-1], sp.Pieces[i]
		h := left.X1 - left.X0
		d1 := left.B + 2*left.C*h + 3*left.D*h*h - right.B
		d2 := 2*left.C + 6*left.D*h - 2*right.C
		fmt.Printf("  x%d = %g: |ΔS'| = %g, |ΔS''| = %g\n", i, right.X0, abs(d1), abs(d2))
	}
	d0, err := sp.SecondDerivative(sp.X[0])
	if err != nil {
		fmt.Println("ошибка:", err)
		return
	}
	dn, err := sp.SecondDerivative(sp.X[len(sp.X)-1])
	if err != nil {
		fmt.Println("ошибка:", err)
		return
	}
	fmt.Printf("граничные условия: S''(x0) = %g, S''(x%d) = %g\n", d0, len(sp.X)-1, dn)
}

func abs(v float64) float64 {
	if v < 0 {
		return -v
	}
	return v
}
