package demo

import (
	"fmt"

	"lab1/input"
	"lab1/iterative"
	"lab1/lu"
)

var maxIter = 10000

func Iterative() {
	fmt.Println("=== 1.3 Простые итерации и метод Зейделя ===")
	a, b, eps, err := input.ReadSystemEps("data/iterative.txt")
	if err != nil {
		fmt.Println("ошибка:", err)
		return
	}
	ref, err := lu.Solve(a, b)
	if err != nil {
		fmt.Println("ошибка:", err)
		return
	}
	fmt.Println("ε =", eps)
	fmt.Println("эталонное решение (LU):")
	fmt.Print(ref)
	xSI, nSI, err := iterative.SimpleIterations(a, b, eps, maxIter)
	if err != nil {
		fmt.Println("ошибка:", err)
		return
	}
	xSeidel, nSeidel, err := iterative.Seidel(a, b, eps, maxIter)
	if err != nil {
		fmt.Println("ошибка:", err)
		return
	}
	fmt.Println("простые итерации:")
	fmt.Print(xSI)
	fmt.Println("итераций (простые итерации):", nSI)
	fmt.Println("метод Зейделя:")
	fmt.Print(xSeidel)
	fmt.Println("итераций (Зейдель):", nSeidel)
	if nSeidel < nSI {
		fmt.Println("Зейдель сошёлся быстрее")
	} else if nSeidel > nSI {
		fmt.Println("простые итерации сошлись быстрее")
	} else {
		fmt.Println("методы сошлись за одинаковое число итераций")
	}
}
