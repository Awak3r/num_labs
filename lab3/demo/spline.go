package demo

import (
	"fmt"

	"lab3/input"
	"lab3/spline"
)

func Spline() {
	fmt.Println("=== 3.2 Кубический сплайн (S''(x0)=S''(x4)=0), x* = 1.5 ===")
	pts, err := input.ReadPointsXstar("data/spline.txt")
	if err != nil {
		fmt.Println("ошибка:", err)
		return
	}
	fmt.Println("исходные узлы (xi, yi):")
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
	fmt.Printf("S(x*) = S_%d(%g) = %12.6f (интервал %d: [%g, %g])\n", idx, pts.XStar, v, idx, sp.Pieces[idx-1].X0, sp.Pieces[idx-1].X1)
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
	fmt.Printf("граничные условия: S''(x0) = %g (должно быть 0), S''(x4) = %g (должно быть 0)\n", d0, dn)
}

func abs(v float64) float64 {
	if v < 0 {
		return -v
	}
	return v
}
