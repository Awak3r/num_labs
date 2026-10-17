package demo

import (
	"fmt"

	"lab3/funcs"
	"lab3/input"
	"lab3/integration"
)

func Integrate() {
	fmt.Println("=== 3.5 Численное интегрирование: ∫ x/(3x+4)³ dx, x0=−1, xk=1 ===")
	x0, xk, h1, h2, err := input.ReadIntegral("data/integral.txt")
	if err != nil {
		fmt.Println("ошибка:", err)
		return
	}
	fmt.Printf("пределы: x0 = %g, xk = %g; шаги: h1 = %g, h2 = %g\n", x0, xk, h1, h2)
	exact := funcs.ExactIntegral(x0, xk)
	fmt.Printf("точное значение: F = %12.8f\n", exact)

	fmt.Printf("%-24s %14s %14s %4s %16s %14s %14s\n", "метод", "J(h1)", "J(h2)", "p", "J_Рунге-Ромберг", "погр. J(h2)", "|J(h2)−F|")
	var rows []struct {
		name    string
		absDiff float64
	}
	for _, m := range integration.Methods {
		n1 := int((xk-x0)/h1 + 0.5)
		n2 := int((xk-x0)/h2 + 0.5)
		j1, err := m.Fn(funcs.Integrand, x0, xk, n1)
		if err != nil {
			fmt.Printf("%-24s ошибка: %v\n", m.Name, err)
			continue
		}
		j2, err := m.Fn(funcs.Integrand, x0, xk, n2)
		if err != nil {
			fmt.Printf("%-24s ошибка: %v\n", m.Name, err)
			continue
		}
		rr, e2 := integration.RungeRomberg(j1, j2, m.P)
		fmt.Printf("%-24s %14.8f %14.8f %4d %16.8f %14.8f %14.8f\n", m.Name, j1, j2, m.P, rr, abs(e2), abs(j2-exact))
		rows = append(rows, struct {
			name    string
			absDiff float64
		}{m.Name, abs(j2 - exact)})
	}
	best := rows[0]
	for _, r := range rows[1:] {
		if r.absDiff < best.absDiff {
			best = r
		}
	}
	fmt.Printf("наиболее точный метод (мин. |J(h2)−F|): %s (%g)\n", best.name, best.absDiff)
}
