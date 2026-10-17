package demo

import (
	"fmt"

	"lab1/input"
	"lab1/matrix"
	"lab1/qr"
)

func QR() {
	fmt.Println("=== 1.5 QR-разложение ===")
	a, err := input.ReadMatrix("data/qr.txt")
	if err != nil {
		fmt.Println("ошибка:", err)
		return
	}
	printMatrix("A", a)

	q, r, err := qr.Factor(a)
	if err != nil {
		fmt.Println("ошибка:", err)
		return
	}
	printMatrix("Q", q)
	printMatrix("R", r)
	qrProd, _ := q.Multiply(r)
	printMatrix("Q·R", qrProd)

	real, pairs, err := qr.Find(a, 1000)
	if err != nil {
		fmt.Println("ошибка:", err)
		return
	}
	if len(real) > 0 {
		fmt.Println("вещественные собственные значения:", real)
	}
	printPairs(pairs)

	fmt.Println("--- свой пример: матрица с комплексно-сопряжёнными СЗ ---")
	ex, err := input.ReadMatrix("data/qr_example.txt")
	if err != nil {
		fmt.Println("ошибка:", err)
		return
	}
	printMatrix("A", ex)
	_, exPairs, err := qr.Find(ex, 1000)
	if err != nil {
		fmt.Println("ошибка:", err)
		return
	}
	if len(exPairs) == 0 {
		fmt.Println("комплексная пара не обнаружена qr.Find — неожиданно")
		return
	}
	printPairs(exPairs)
}

func printPairs(pairs []qr.Pair) {
	if len(pairs) == 0 {
		return
	}
	fmt.Println("комплексно-сопряжённая пара:")
	for _, p := range pairs {
		fmt.Printf("λ = %.6f %+.6fi\n", p.Re, p.Im)
	}
}

func printMatrix(name string, m *matrix.Matrix) {
	fmt.Println(name, "=")
	fmt.Print(m)
}
