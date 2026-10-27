package interpolation

import (
	"fmt"
	"strings"
)

type Node struct {
	X, Y float64
}

type Result struct {
	Coef []float64
	Form string
}

func DividedDifferences(nodes []Node) [][]float64 {
	n := len(nodes)
	d := make([][]float64, n)
	d[0] = make([]float64, n)
	for i := range nodes {
		d[0][i] = nodes[i].Y
	}
	for k := 1; k < n; k++ {
		d[k] = make([]float64, n-k)
		for i := 0; i < n-k; i++ {
			d[k][i] = (d[k-1][i+1] - d[k-1][i]) / (nodes[i+k].X - nodes[i].X)
		}
	}
	return d
}

func Newton(nodes []Node) Result {
	n := len(nodes)
	div := DividedDifferences(nodes)
	coefs := make([]float64, n)
	for k := 0; k < n; k++ {
		coefs[k] = div[k][0]
	}
	var sb strings.Builder
	sb.WriteString("N(x) = ")
	var terms []string
	prefix := ""
	for k := 0; k < n; k++ {
		if k == 0 {
			terms = append(terms, fmt.Sprintf("%.*g", 6, coefs[k]))
		} else {
			if k > 1 {
				prefix += fmt.Sprintf("(x%+.*g)", 4, -nodes[k-1].X)
			} else {
				prefix = fmt.Sprintf("(x%+.*g)", 4, -nodes[0].X)
			}
			terms = append(terms, fmt.Sprintf("%+.*g%s", 6, coefs[k], prefix))
		}
	}
	sb.WriteString(strings.Join(terms, " "))
	return Result{Coef: coefs, Form: sb.String()}
}

func EvalNewton(nodes []Node, coefs []float64, x float64) float64 {
	n := len(coefs)
	res := coefs[n-1]
	for k := n - 2; k >= 0; k-- {
		res = res*(x-nodes[k].X) + coefs[k]
	}
	return res
}

func Lagrange(nodes []Node) Result {
	n := len(nodes)
	coef := make([]float64, n)
	for i := 0; i < n; i++ {
		basis := []float64{1}
		scale := 1.0
		for j := 0; j < n; j++ {
			if j != i {
				basis = polyMulShift(basis, nodes[j].X)
				scale *= nodes[i].X - nodes[j].X
			}
		}
		c := nodes[i].Y / scale
		for k := range basis {
			coef[k] += c * basis[k]
		}
	}
	form := "L(x) = " + polyString(coef)
	return Result{Coef: coef, Form: form}
}

func polyString(coef []float64) string {
	var parts []string
	for i := len(coef) - 1; i >= 0; i-- {
		c := coef[i]
		switch {
		case i == 0:
			parts = append(parts, fmt.Sprintf("%+.*g", 6, c))
		case i == 1:
			parts = append(parts, fmt.Sprintf("%+.*g·x", 6, c))
		default:
			parts = append(parts, fmt.Sprintf("%+.*g·x^%d", 6, c, i))
		}
	}
	return strings.Join(parts, " ")
}

func EvalPoly(coef []float64, x float64) float64 {
	res := 0.0
	for i := len(coef) - 1; i >= 0; i-- {
		res = res*x + coef[i]
	}
	return res
}

func polyMulShift(a []float64, root float64) []float64 {
	r := make([]float64, len(a)+1)
	for i, c := range a {
		r[i+1] += c
		r[i] += -root * c
	}
	return r
}
