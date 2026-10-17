package demo

import (
	"fmt"

	"lab3/differentiation"
	"lab3/input"
)

func Derivative() {
	fmt.Println("=== 3.4 Численное дифференцирование, x* = 2.0 ===")
	pts, err := input.ReadPointsXstar("data/derivative.txt")
	if err != nil {
		fmt.Println("ошибка:", err)
		return
	}
	fmt.Println("исходная таблица значений функции:")
	for i := range pts.X {
		fmt.Printf("  x%d = %8.4f   y%d = %12.6f\n", i, pts.X[i], i, pts.Y[i])
	}
	fmt.Printf("целевая точка: x* = %g\n", pts.XStar)
	res, err := differentiation.Compute(pts.X, pts.Y, pts.XStar)
	if err != nil {
		fmt.Println("ошибка:", err)
		return
	}
	fmt.Print("узлы, использованные в расчёте: ")
	for _, i := range res.Used {
		fmt.Printf("(x%d=%g, y%d=%g) ", i, pts.X[i], i, pts.Y[i])
	}
	fmt.Println()
	fmt.Printf("первая производная:  y'(x*) = %12.8f  [%s]\n", res.D1, res.D1Form)
	fmt.Printf("вторая производная: y''(x*) = %12.8f  [%s]\n", res.D2, res.D2Form)
}
