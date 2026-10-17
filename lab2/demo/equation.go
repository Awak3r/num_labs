package demo

import (
	"fmt"

	"lab2/equation"
	"lab2/funcs"
	"lab2/input"
)

func Equation() {
	fmt.Println("=== 2.1 Нелинейное уравнение: sqrt(1-x^2) - e^x + 0.1 = 0 ===")
	eps, x0, err := input.ReadEquation("data/equation.txt")
	if err != nil {
		fmt.Println("ошибка:", err)
		return
	}
	fmt.Println("ε =", eps)
	fmt.Println("начальное приближение (графически): x0 =", x0)
	fmt.Println("эталон (сборник): положительный корень x ≈ 0.0915")

	xSI, nSI, stepsSI, err := equation.SimpleIterations(funcs.Phi, funcs.F, x0, eps, maxIter)
	if err != nil {
		fmt.Println("ошибка:", err)
		return
	}
	xN, nN, stepsN, err := equation.Newton(funcs.F, funcs.DF, x0, eps, maxIter)
	if err != nil {
		fmt.Println("ошибка:", err)
		return
	}

	fmt.Println("итерация |      Δ SI          x_SI     |    Δ Ньютон        x_Ньютон")
	k := nSI
	if nN > k {
		k = nN
	}
	get := func(steps []equation.Step, i int) (float64, float64, bool) {
		if i < len(steps) {
			return steps[i].Diff, steps[i].X, true
		}
		return 0, 0, false
	}
	for i := 0; i < k; i++ {
		d1, x1, ok1 := get(stepsSI, i)
		d2, x2, ok2 := get(stepsN, i)
		c1, c2 := "     -            -", "     -            -"
		if ok1 {
			c1 = fmt.Sprintf("%.3e  %12.8f", d1, x1)
		}
		if ok2 {
			c2 = fmt.Sprintf("%.3e  %12.8f", d2, x2)
		}
		fmt.Printf("   %3d   |  %s  |  %s\n", i+1, c1, c2)
	}

	fmt.Printf("корень (простые итерации): x = %.6f, итераций: %d\n", xSI, nSI)
	fmt.Printf("корень (Ньютон):           x = %.6f, итераций: %d\n", xN, nN)
}
