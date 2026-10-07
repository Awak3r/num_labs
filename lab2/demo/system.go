package demo

import (
	"fmt"

	"lab2/funcs"
	"lab2/input"
	"lab2/matrix"
	"lab2/system"
)

func System() {
	fmt.Println("=== 2.1 Система НУ (a = 4): (x1^2+16)x2-64=0; (x1-2)^2+(x2-2)^2-16=0 ===")
	_, eps, x0v, err := input.ReadSystem("data/system.txt")
	if err != nil {
		fmt.Println("ошибка:", err)
		return
	}
	x0 := matrix.Column(x0v)
	fmt.Println("ε =", eps)
	fmt.Println("начальное приближение (графически): x0 =", funcs.VectorString(x0))
	fmt.Println("эталон (сборник): x1 = 5.9293, x2 = 1.2511")

	fmt.Println("--- метод простой итерации ---")
	xSI, nSI, stepsSI, err := system.SimpleIterations(funcs.PhiSystem, x0, eps, maxIter)
	if err != nil {
		fmt.Println("ошибка:", err)
	} else {
		printStepsVector(stepsSI)
		fmt.Printf("решение (простые итерации): %s, итераций: %d\n", funcs.VectorString(xSI), nSI)
		fmt.Printf("невязка: f1 = %+.2e, f2 = %+.2e\n", funcs.F1(xSI), funcs.F2(xSI))
	}

	fmt.Println("--- метод Ньютона ---")
	xN, nN, stepsN, err := system.Newton(funcs.FSystem, funcs.Jacobian, x0, eps, maxIter)
	if err != nil {
		fmt.Println("ошибка:", err)
		return
	}
	printStepsVector(stepsN)
	fmt.Printf("решение (Ньютон): %s, итераций: %d\n", funcs.VectorString(xN), nN)
	fmt.Printf("невязка: f1 = %+.2e, f2 = %+.2e\n", funcs.F1(xN), funcs.F2(xN))

	fmt.Println("--- анализ сходимости ---")
	fmt.Println("итерация |  |Δ| простые итерации  |  |Δ| Ньютон")
	k := nSI
	if nN > k {
		k = nN
	}
	get := func(steps []system.Step, i int) (float64, bool) {
		if i < len(steps) {
			return steps[i].Diff, true
		}
		return 0, false
	}
	for i := 0; i < k; i++ {
		s1, s2 := "-", "-"
		if d, ok := get(stepsSI, i); ok {
			s1 = fmt.Sprintf("%.3e", d)
		}
		if d, ok := get(stepsN, i); ok {
			s2 = fmt.Sprintf("%.3e", d)
		}
		fmt.Printf("   %3d   |  %-20s  |  %s\n", i+1, s1, s2)
	}
	fmt.Printf("Ньютон сошёлся за %d итераций, простые итерации — за %d: ", nN, nSI)
	if nN < nSI {
		fmt.Println("Ньютон сходится быстрее (квадратичная сходимость)")
	} else if nN > nSI {
		fmt.Println("простые итерации сошлись быстрее")
	} else {
		fmt.Println("методы сошлись за одинаковое число итераций")
	}
}

func printStepsVector(steps []system.Step) {
	fmt.Println(" k  |     x1_k        |     x2_k        |   |Δ|")
	for _, s := range steps {
		fmt.Printf("%2d  |  %14.8f  |  %14.8f  |  %.3e\n", s.Iter, s.X.Get(0, 0), s.X.Get(1, 0), s.Diff)
	}
}
