package demo

import (
	"fmt"
	"os/exec"

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

	real, hasPair, err := qr.Find(a, 1000)
	if err != nil {
		fmt.Println("ошибка:", err)
		return
	}
	if len(real) > 0 {
		fmt.Println("вещественные собственные значения:", real)
	}
	if hasPair {
		fmt.Println("комплексно-сопряжённая пара:")
		out, err := exec.Command("python3", "qr/pair.py").Output()
		if err != nil {
			fmt.Println("не удалось запустить qr/pair.py (нужен python3+NumPy):", err)
			return
		}
		fmt.Print(string(out))
	}

	fmt.Println("--- свой пример: матрица с комплексно-сопряжёнными СЗ ---")
	ex, err := input.ReadMatrix("data/qr_example.txt")
	if err != nil {
		fmt.Println("ошибка:", err)
		return
	}
	printMatrix("A", ex)
	_, exHasPair, err := qr.Find(ex, 1000)
	if err != nil {
		fmt.Println("ошибка:", err)
		return
	}
	if !exHasPair {
		fmt.Println("комплексная пара не обнаружена qr.Find — неожиданно")
		return
	}
	fmt.Println("комплексно-сопряжённая пара:")
	out, err := exec.Command("python3", "qr/pair.py", "data/qr_example.txt").Output()
	if err != nil {
		fmt.Println("не удалось запустить qr/pair.py (нужен python3+NumPy):", err)
		return
	}
	fmt.Print(string(out))
}

func printMatrix(name string, m *matrix.Matrix) {
	fmt.Println(name, "=")
	fmt.Print(m)
}
