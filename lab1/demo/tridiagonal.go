package demo

import (
	"fmt"

	"lab1/input"
	"lab1/tridiagonal"
)

func Tridiagonal() {
	fmt.Println("=== 1.2 Метод прогонки ===")
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
	fmt.Println("прогоночные α =", alpha)
	fmt.Println("прогоночные β =", beta)
	fmt.Println("x =")
	fmt.Print(x)
}
