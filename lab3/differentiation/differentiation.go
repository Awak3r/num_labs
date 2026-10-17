package differentiation

import (
	"fmt"
	"math"
)

// Result — результат вычисления производных с описанием использованных узлов.
type Result struct {
	Used   []int   // индексы узлов, использованных в расчёте
	D1     float64 // первая производная
	D2     float64 // вторая производная
	D1Form string
	D2Form string
}

// Compute вычисляет y'(x*) и y''(x*) по таблично заданной функции.
// Если x* — внутренний узел (не на краю), используются центральные формулы O(h²):
//
//	y' ≈ (y[i+1] − y[i−1]) / (2h)
//	y'' ≈ (y[i+1] − 2y[i] + y[i−1]) / h²
//
// На краях — односторонние формулы второго порядка.
func Compute(x, y []float64, xstar float64) (Result, error) {
	n := len(x)
	if n < 4 || len(y) != n {
		return Result{}, fmt.Errorf("differentiation: нужно ≥4 узлов, есть %d", n)
	}
	// ищем узел, совпадающий с x* (в пределах допуска)
	idx := -1
	for i := range x {
		if math.Abs(x[i]-xstar) < 1e-9 {
			idx = i
			break
		}
	}
	if idx < 0 {
		return Result{}, fmt.Errorf("differentiation: x*=%g не совпадает ни с одним узлом", xstar)
	}
	h := x[1] - x[0]
	for i := 1; i < n; i++ {
		if math.Abs(x[i]-x[i-1]-h) > 1e-9 {
			return Result{}, fmt.Errorf("differentiation: сетка неравномерная — формулы требуют равного шага")
		}
	}
	res := Result{}
	switch {
	case idx >= 2 && idx <= n-3:
		// пяти-точечная центральная формула O(h⁴) для первой производной
		res.D1 = (y[idx-2] - 8*y[idx-1] + 8*y[idx+1] - y[idx+2]) / (12 * h)
		res.Used = []int{idx - 2, idx - 1, idx, idx + 1, idx + 2}
		res.D1Form = "(y[i−2] − 8y[i−1] + 8y[i+1] − y[i+2]) / (12h)"
		res.D2 = (-y[idx-2] + 16*y[idx-1] - 30*y[idx] + 16*y[idx+1] - y[idx+2]) / (12 * h * h)
		res.D2Form = "(−y[i−2] + 16y[i−1] − 30y[i] + 16y[i+1] − y[i+2]) / (12h²)"
	case idx == 1 || idx == n-2:
		// центральная 3-точечная O(h²)
		res.D1 = (y[idx+1] - y[idx-1]) / (2 * h)
		res.Used = []int{idx - 1, idx, idx + 1}
		res.D1Form = "(y[i+1] − y[i−1]) / (2h)"
		res.D2 = (y[idx+1] - 2*y[idx] + y[idx-1]) / (h * h)
		res.D2Form = "(y[i+1] − 2y[i] + y[i−1]) / h²"
	default:
		// края: односторонние O(h²)
		if idx == 0 {
			res.D1 = (-3*y[0] + 4*y[1] - y[2]) / (2 * h)
			res.Used = []int{0, 1, 2}
			res.D1Form = "(−3y[0] + 4y[1] − y[2]) / (2h)"
			res.D2 = (2*y[0] - 5*y[1] + 4*y[2] - y[3]) / (h * h)
			res.D2Form = "(2y[0] − 5y[1] + 4y[2] − y[3]) / h²"
		} else {
			i := n - 1
			res.D1 = (3*y[i] - 4*y[i-1] + y[i-2]) / (2 * h)
			res.Used = []int{i - 2, i - 1, i}
			res.D1Form = "(3y[i] − 4y[i−1] + y[i−2]) / (2h)"
			res.D2 = (2*y[i] - 5*y[i-1] + 4*y[i-2] - y[i-3]) / (h * h)
			res.D2Form = "(2y[i] − 5y[i−1] + 4y[i−2] − y[i−3]) / h²"
		}
	}
	return res, nil
}
