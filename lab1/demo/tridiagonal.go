package demo

import (
	"fmt"

	"lab1/input"
	"lab1/matrix"
	"lab1/tridiagonal"
)

func Tridiagonal() {
	main, super, sub, d, err := input.ReadTridiagonal("data/tridiagonal.txt")
	if err != nil {
		fmt.Println("ошибка:", err)
		return
	}

	x, alpha, beta, err := tridiagonal.Solve(main, super, sub, d)
	if err != nil {
		fmt.Println("ошибка:", err)
		return
	}

	fmt.Println("главная диагональ =", main)
	fmt.Println("верхняя диагональ =", super)
	fmt.Println("нижняя диагональ  =", sub)
	fmt.Println("правая часть d    =", d)
	fmt.Println("прогоночные α =", alpha)
	fmt.Println("прогоночные β =", beta)
	fmt.Println("x =")
	fmt.Print(x)

	a := buildTridiagonal(main, super, sub)
	ax, _ := a.Multiply(x)
	fmt.Println("A·x")
	fmt.Print(ax)
	fmt.Println("A·x == d:", approxEqual(ax, matrix.Column(d)))
}

func buildTridiagonal(main, super, sub []float64) *matrix.Matrix {
	n := len(main)
	m := matrix.Zeros(n, n)
	for i := 0; i < n; i++ {
		m.Set(i, i, main[i])
	}
	for i := 0; i < n-1; i++ {
		m.Set(i, i+1, super[i])
	}
	for i := 1; i < n; i++ {
		m.Set(i, i-1, sub[i-1])
	}
	return m
}
