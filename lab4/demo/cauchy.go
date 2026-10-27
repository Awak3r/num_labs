package demo

import (
	"fmt"

	"lab4/funcs"
	"lab4/input"
	"lab4/ode"
)

func Cauchy() {
	fmt.Println("=== 4.1 Задача Коши: y'' = 2y + 4x^2·e^(x^2), y(0) = 3, y'(0) = 0 ===")
	x0, xk, h1, h2, err := input.ReadInterval("data/cauchy.txt")
	if err != nil {
		fmt.Println("ошибка:", err)
		return
	}
	fmt.Printf("отрезок: x ∈ [%g, %g]; шаги: h1 = %g, h2 = %g\n", x0, xk, h1, h2)
	fmt.Println("точное решение: y(x) = e^(x^2) + e^(x√2) + e^(−x√2)")

	xsE, usE, err := ode.Euler(funcs.FCauchy, x0, xk, h1, funcs.CauchyU0())
	if err != nil {
		fmt.Println("ошибка:", err)
		return
	}
	xsR, usR, err := ode.RK4(funcs.FCauchy, x0, xk, h1, funcs.CauchyU0())
	if err != nil {
		fmt.Println("ошибка:", err)
		return
	}
	xsA, usA, err := ode.Adams4(funcs.FCauchy, x0, xk, h1, funcs.CauchyU0())
	if err != nil {
		fmt.Println("ошибка:", err)
		return
	}
	xsE2, usE2, err := ode.Euler(funcs.FCauchy, x0, xk, h2, funcs.CauchyU0())
	if err != nil {
		fmt.Println("ошибка:", err)
		return
	}
	xsR2, usR2, err := ode.RK4(funcs.FCauchy, x0, xk, h2, funcs.CauchyU0())
	if err != nil {
		fmt.Println("ошибка:", err)
		return
	}
	xsA2, usA2, err := ode.Adams4(funcs.FCauchy, x0, xk, h2, funcs.CauchyU0())
	if err != nil {
		fmt.Println("ошибка:", err)
		return
	}

	ysE, ysR, ysA := col(usE), col(usR), col(usA)
	ysE2, ysR2, ysA2 := col(usE2), col(usR2), col(usA2)

	fmt.Printf("численное решение с шагом h1 = %g:\n", h1)
	fmt.Printf("%6s %12s %12s %12s %12s\n", "x", "Эйлер", "РК4", "Адамс-4", "точное")
	for i := range xsE {
		fmt.Printf("%6.2f %12.6f %12.6f %12.6f %12.6f\n",
			xsE[i], ysE[i], ysR[i], ysA[i], funcs.ExactCauchy(xsE[i]))
	}

	fmt.Println("погрешности численного решения (макс. по узлам):")
	fmt.Printf("%-10s %3s %20s %22s\n", "метод", "p", "ε_Рунге-Ромберг", "ε_сравнение_с_точным")
	rows := []struct {
		name string
		p    int
		rr   float64
		ex   float64
	}{
		{"Эйлер", 1, 0, 0},
		{"РК4", 4, 0, 0},
		{"Адамс-4", 4, 0, 0},
	}
	for j, r := range []struct {
		xs1, ys1, xs2, ys2 []float64
	}{
		{xsE, ysE, xsE2, ysE2},
		{xsR, ysR, xsR2, ysR2},
		{xsA, ysA, xsA2, ysA2},
	} {
		rows[j].rr, err = rungeError(h1, h2, rows[j].p, r.xs1, r.ys1, r.xs2, r.ys2)
		if err != nil {
			fmt.Println("ошибка:", err)
			return
		}
		rows[j].ex = maxAbsDiff(r.xs1, r.ys1, funcs.ExactCauchy)
	}
	for _, r := range rows {
		fmt.Printf("%-10s %3d %20g %22g\n", r.name, r.p, r.rr, r.ex)
	}
}
