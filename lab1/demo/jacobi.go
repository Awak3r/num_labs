package demo

import (
	"fmt"

	"lab1/input"
	"lab1/jacobi"
	"lab1/matrix"
)

func Jacobi() {
	fmt.Println("=== 1.4 Вращения Якоби ===")
	a, eps, err := input.ReadMatrixEps("data/jacobi.txt")
	if err != nil {
		fmt.Println("ошибка:", err)
		return
	}

	fmt.Println("A =")
	fmt.Print(a)
	fmt.Println("ε =", eps)

	vals, v, iters, err := jacobi.Solve(a, eps, 10000)
	if err != nil {
		fmt.Println("ошибка:", err)
		return
	}
	fmt.Println("собственные значения:", vals)
	fmt.Println("матрица собственных векторов V =")
	fmt.Print(v)
	fmt.Println("итераций (вращений):", iters)

	lam := matrix.Diagonal(vals)
	av, _ := a.Multiply(v)
	vlam, _ := v.Multiply(lam)
	fmt.Println("A·V =")
	fmt.Print(av)
	fmt.Println("V·Λ =")
	fmt.Print(vlam)
	fmt.Println("A·V == V·Λ (допуск ε):", av.ApproxEqualTol(vlam, eps))
}
