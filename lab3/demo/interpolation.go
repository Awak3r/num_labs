package demo

import (
	"fmt"
	"math"

	"lab3/funcs"
	"lab3/input"
	"lab3/interpolation"
)

func printDividedDiff(nodes []interpolation.Node) {
	d := interpolation.DividedDifferences(nodes)
	fmt.Println("таблица разделённых разностей (k — порядок разности):")
	fmt.Printf("      ")
	for i := range nodes {
		fmt.Printf("%14.4f", nodes[i].X)
	}
	fmt.Println()
	for k, row := range d {
		fmt.Printf("k=%d: ", k)
		for _, v := range row {
			fmt.Printf("%14.6f", v)
		}
		fmt.Println()
	}
}

func interpolationSet(name, path string) {
	fmt.Printf("--- набор %s (%s) ---\n", name, path)
	xs, xstar, err := input.ReadNodesXstar(path)
	if err != nil {
		fmt.Println("ошибка:", err)
		return
	}
	nodes := make([]interpolation.Node, len(xs))
	for i, xv := range xs {
		nodes[i] = interpolation.Node{X: xv, Y: funcs.Tan(xv)}
	}
	fmt.Println("исходные узлы (xi, yi=tan(xi)):")
	for i := range nodes {
		fmt.Printf("  x%d = %12.8f   y%d = %12.8f\n", i, nodes[i].X, i, nodes[i].Y)
	}
	fmt.Printf("x* = %g\n", xstar)

	lag := interpolation.Lagrange(nodes)
	newt := interpolation.Newton(nodes)
	fmt.Println("многочлен Лагранжа (аналитический вид):")
	fmt.Println("  " + lag.Form)
	printDividedDiff(nodes)
	fmt.Println("многочлен Ньютона (аналитический вид):")
	fmt.Println("  " + newt.Form)

	lv := interpolation.EvalPoly(lag.Coef, xstar)
	nv := interpolation.EvalNewton(nodes, newt.Coef, xstar)
	fv := funcs.Tan(xstar)
	fmt.Printf("L(x*) = %12.8f\n", lv)
	fmt.Printf("N(x*) = %12.8f\n", nv)
	fmt.Printf("f(x*) = tan(x*) = %12.8f\n", fv)
	fmt.Printf("абсолютная погрешность интерполяции |f(x*)−L(x*)| = %g\n", math.Abs(fv-lv))
	fmt.Printf("|L(x*)−N(x*)| = %g\n", math.Abs(lv-nv))
}

func Interpolation() {
	fmt.Println("=== 3.1 Интерполяция: Лагранж и Ньютон, y = tg(x), x* = 3π/16 ===")
	interpolationSet("а", "data/interpolation_a.txt")
	fmt.Println()
	interpolationSet("б", "data/interpolation_b.txt")
}
