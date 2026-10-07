package demo

import (
	"fmt"

	"lab1/input"
	"lab1/lu"
	"lab1/matrix"
)

func LU() {
	fmt.Println("=== 1.1 LU-разложение ===")
	a, b, err := input.ReadSystem("data/lu.txt")
	if err != nil {
		fmt.Println("ошибка:", err)
		return
	}

	l, u, p, err := lu.Decompose(a)
	if err != nil {
		fmt.Println("ошибка:", err)
		return
	}
	luProd, _ := l.Multiply(u)
	pa, _ := p.Multiply(a)

	fmt.Println("A =")
	fmt.Print(a)
	fmt.Println("b =")
	fmt.Print(b)
	fmt.Println("L =")
	fmt.Print(l)
	fmt.Println("U =")
	fmt.Print(u)
	fmt.Println("P =")
	fmt.Print(p)
	fmt.Println("L·U =")
	fmt.Print(luProd)
	fmt.Println("P·A =")
	fmt.Print(pa)
	fmt.Println("L·U == P·A:", luProd.ApproxEqual(pa))

	x, err := lu.Solve(a, b)
	if err != nil {
		fmt.Println("ошибка:", err)
		return
	}
	fmt.Println("x =")
	fmt.Print(x)
	ax, _ := a.Multiply(x)
	fmt.Println("A·x == b:", ax.ApproxEqual(b))

	inv, err := lu.Inverse(a)
	if err != nil {
		fmt.Println("ошибка:", err)
		return
	}
	fmt.Println("A⁻¹ =")
	fmt.Print(inv)
	ainv, _ := a.Multiply(inv)
	fmt.Println("A·A⁻¹ =")
	fmt.Print(ainv)
	fmt.Println("A·A⁻¹ == E:", ainv.ApproxEqual(matrix.MakeIdentityMatrix(a.Rows())))

	det, err := lu.Det(a)
	if err != nil {
		fmt.Println("ошибка:", err)
		return
	}
	fmt.Println("det(A) =", det)
}
