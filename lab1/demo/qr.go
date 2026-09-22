package demo

import (
	"fmt"
	"os/exec"

	"lab1/input"
	"lab1/lu"
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
	fmt.Println("Q·R == A:", isEqual(qrProd, a))

	real, hasPair, err := qr.Eigenvalues(a, 1000)
	if err != nil {
		fmt.Println("ошибка:", err)
		return
	}
	if len(real) > 0 {
		fmt.Println("вещественные собственные значения:", real)
	}
	if hasPair {
		fmt.Println("комплексно-сопряжённая пара (через NumPy, разрешено ТЗ):")
		out, err := exec.Command("python3", "qr/pair.py").Output()
		if err != nil {
			fmt.Println("не удалось запустить qr/pair.py (нужен python3+NumPy):", err)
			return
		}
		fmt.Print(string(out))
		return
	}

	tr := 0.0
	for i := 0; i < a.Rows(); i++ {
		tr += a.Get(i, i)
	}
	sum, prod := 0.0, 1.0
	for _, l := range real {
		sum += l
		prod *= l
	}
	det, err := lu.Det(a)
	if err != nil {
		fmt.Println("ошибка:", err)
		return
	}
	fmt.Printf("проверка: сумма СЗ = %g, trace(A) = %g\n", sum, tr)
	fmt.Printf("проверка: произведение СЗ = %g, det(A) = %g\n", prod, det)
}

func printMatrix(name string, m *matrix.Matrix) {
	fmt.Println(name, "=")
	fmt.Print(m)
}
