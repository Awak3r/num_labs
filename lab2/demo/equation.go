package demo

import (
	"fmt"
	"math"

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

	fmt.Println("--- проверка условия сходимости: |phi'(x)| < 1 ---")
	for _, x := range []float64{x0, 0.05, 0.09, 0.2} {
		fmt.Printf("phi'(%g) = %+.4f\n", x, funcs.DPhi(x))
	}

	fmt.Println("--- метод простой итерации: x = ln(sqrt(1-x^2)+0.1) ---")
	xSI, nSI, stepsSI, err := equation.SimpleIterations(funcs.Phi, funcs.F, x0, eps, maxIter)
	if err != nil {
		fmt.Println("ошибка:", err)
		return
	}
	printStepsScalar(stepsSI)
	fmt.Printf("корень (простые итерации): x = %.6f, f(x) = %+.2e, итераций: %d\n", xSI, funcs.F(xSI), nSI)

	fmt.Println("--- метод Ньютона ---")
	xN, nN, stepsN, err := equation.Newton(funcs.F, funcs.DF, x0, eps, maxIter)
	if err != nil {
		fmt.Println("ошибка:", err)
		return
	}
	printStepsScalar(stepsN)
	fmt.Printf("корень (Ньютон): x = %.6f, f(x) = %+.2e, итераций: %d\n", xN, funcs.F(xN), nN)

	fmt.Println("--- анализ сходимости ---")
	fmt.Println("итерация |  погрешность SI  |  погрешность Ньютона")
	k := nSI
	if nN > k {
		k = nN
	}
	get := func(steps []equation.Step, i int) float64 {
		if i < len(steps) {
			return steps[i].Diff
		}
		return math.NaN()
	}
	for i := 0; i < k; i++ {
		d1, d2 := get(stepsSI, i), get(stepsN, i)
		s1, s2 := "-", "-"
		if !math.IsNaN(d1) {
			s1 = fmt.Sprintf("%.3e", d1)
		}
		if !math.IsNaN(d2) {
			s2 = fmt.Sprintf("%.3e", d2)
		}
		fmt.Printf("   %3d   |  %-14s  |  %s\n", i+1, s1, s2)
	}
	fmt.Printf("Ньютон сошёлся за %d итераций, простые итерации — за %d: ", nN, nSI)
	if nN < nSI {
		fmt.Println("Ньютон сходится быстрее (квадратичная сходимость против линейной у простых итераций)")
	} else if nN > nSI {
		fmt.Println("простые итерации сошлись быстрее")
	} else {
		fmt.Println("методы сошлись за одинаковое число итераций")
	}
}

func printStepsScalar(steps []equation.Step) {
	fmt.Println(" k  |    x_k       |   f(x_k)      |  |x_k - x_{k-1}|")
	for _, s := range steps {
		fmt.Printf("%2d  |  %12.8f  |  %+.2e  |  %.3e\n", s.Iter, s.X, s.Fx, s.Diff)
	}
}
