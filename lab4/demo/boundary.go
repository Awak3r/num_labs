package demo

import (
	"fmt"

	"lab4/finitediff"
	"lab4/funcs"
	"lab4/input"
	"lab4/shooting"
)

func Boundary() {
	fmt.Println("=== 4.2 Краевая задача: x^2·(x+1)·y'' = 2y, y'(1) = −1, 2y(2) − 4y'(2) = 4 ===")
	x0, xk, h1, h2, eps, err := input.ReadBoundary("data/boundary.txt")
	if err != nil {
		fmt.Println("ошибка:", err)
		return
	}
	fmt.Printf("отрезок: x ∈ [%g, %g]; шаги: h1 = %g, h2 = %g; eps = %g\n", x0, xk, h1, h2, eps)
	fmt.Println("точное решение: y(x) = 1 + 1/x")

	// --- метод стрельбы: подбирается s = y(x0) ---
	s, iters, xsS, usS, err := shooting.Solve(funcs.FBoundary, x0, xk, funcs.DU0(),
		funcs.ShootS0(), funcs.ShootS1(), h1, eps, 100, funcs.RightResidual)
	if err != nil {
		fmt.Println("ошибка:", err)
		return
	}
	fmt.Println("--- метод стрельбы ---")
	fmt.Printf("s = y(%g) = %g, итераций секущих: %d\n", x0, s, iters)
	_, _, xsS2, usS2, err := shooting.Solve(funcs.FBoundary, x0, xk, funcs.DU0(),
		funcs.ShootS0(), funcs.ShootS1(), h2, eps, 100, funcs.RightResidual)
	if err != nil {
		fmt.Println("ошибка:", err)
		return
	}
	ysS, ysS2 := col(usS), col(usS2)

	// --- конечно-разностный метод ---
	a0, b0, g0 := funcs.BC0()
	a1, b1, g1 := funcs.BC1()
	xsD, ysD, err := finitediff.Solve(funcs.PBoundary, funcs.QBoundary, funcs.RBoundary,
		x0, xk, h1, a0, b0, g0, a1, b1, g1)
	if err != nil {
		fmt.Println("ошибка:", err)
		return
	}
	xsD2, ysD2, err := finitediff.Solve(funcs.PBoundary, funcs.QBoundary, funcs.RBoundary,
		x0, xk, h2, a0, b0, g0, a1, b1, g1)
	if err != nil {
		fmt.Println("ошибка:", err)
		return
	}

	fmt.Printf("численное решение с шагом h1 = %g:\n", h1)
	fmt.Printf("%6s %12s %12s %12s\n", "x", "стрельба", "кон.разн.", "точное")
	for i := range xsS {
		fmt.Printf("%6.2f %12.6f %12.6f %12.6f\n",
			xsS[i], ysS[i], ysD[i], funcs.ExactBoundary(xsS[i]))
	}

	fmt.Println("погрешности численного решения (макс. по узлам):")
	fmt.Printf("%-12s %3s %20s %22s\n", "метод", "p", "ε_Рунге-Ромберг", "ε_сравнение_с_точным")
	rrS, err := rungeError(h1, h2, 4, xsS, ysS, xsS2, ysS2)
	if err != nil {
		fmt.Println("ошибка:", err)
		return
	}
	rrD, err := rungeError(h1, h2, 2, xsD, ysD, xsD2, ysD2)
	if err != nil {
		fmt.Println("ошибка:", err)
		return
	}
	fmt.Printf("%-12s %3d %20g %22g\n", "стрельба", 4, rrS, maxAbsDiff(xsS, ysS, funcs.ExactBoundary))
	fmt.Printf("%-12s %3d %20g %22g\n", "кон.разн.", 2, rrD, maxAbsDiff(xsD, ysD, funcs.ExactBoundary))
}
