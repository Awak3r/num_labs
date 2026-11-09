package integration

import (
	"fmt"
	"math"
)

type Method struct {
	Name string
	P    int
	Fn   func(f func(float64) float64, a, b float64, n int) (float64, error)
}

var Methods = []Method{
	{"левые прямоугольники", 1, LeftRect},
	{"правые прямоугольники", 1, RightRect},
	{"средние прямоугольники", 2, MidRect},
	{"трапеции", 2, Trapezoid},
	{"Симпсон", 4, Simpson},
}

func LeftRect(f func(float64) float64, a, b float64, n int) (float64, error) {
	return rects(f, a, b, n, 0)
}

func RightRect(f func(float64) float64, a, b float64, n int) (float64, error) {
	return rects(f, a, b, n, 1)
}

func MidRect(f func(float64) float64, a, b float64, n int) (float64, error) {
	return rects(f, a, b, n, 2)
}

func rects(f func(float64) float64, a, b float64, n, kind int) (float64, error) {
	if n < 1 {
		return 0, fmt.Errorf("integration: n должно быть ≥1, есть %d", n)
	}
	h := (b - a) / float64(n)
	sum := 0.0
	for i := 0; i < n; i++ {
		var xi float64
		switch kind {
		case 0:
			xi = a + float64(i)*h
		case 1:
			xi = a + float64(i+1)*h
		default:
			xi = a + (float64(i)+0.5)*h
		}
		sum += f(xi)
	}
	return sum * h, nil
}

func Trapezoid(f func(float64) float64, a, b float64, n int) (float64, error) {
	if n < 1 {
		return 0, fmt.Errorf("integration: n должно быть ≥1, есть %d", n)
	}
	h := (b - a) / float64(n)
	sum := (f(a) + f(b)) / 2
	for i := 1; i < n; i++ {
		sum += f(a + float64(i)*h)
	}
	return sum * h, nil
}

func Simpson(f func(float64) float64, a, b float64, n int) (float64, error) {
	if n < 2 || n%2 != 0 {
		return 0, fmt.Errorf("integration: Симпсон требует чётного n ≥2, есть %d", n)
	}
	h := (b - a) / float64(n)
	sum := f(a) + f(b)
	for i := 1; i < n; i++ {
		c := 2.0
		if i%2 == 1 {
			c = 4
		}
		sum += c * f(a+float64(i)*h)
	}
	return sum * h / 3, nil
}

// Steps проверяет, что равномерная сетка с заданным шагом покрывает весь отрезок.
func Steps(a, b, h float64) (int, error) {
	if !finite(a) || !finite(b) || !finite(h) || b <= a || h <= 0 {
		return 0, fmt.Errorf("integration: нужны конечные a < b и h > 0")
	}
	steps := (b - a) / h
	n := math.Round(steps)
	if !finite(steps) || n < 1 || n >= float64(int(^uint(0)>>1)) || math.Abs(steps-n) > 1e-9 {
		return 0, fmt.Errorf("integration: шаг %g не делит отрезок [%g, %g] или сетка слишком велика", h, a, b)
	}
	return int(n), nil
}

func RungeRomberg(j1, j2, h1, h2 float64, p int) (refined, err2 float64, err error) {
	if !finite(j1) || !finite(j2) || !finite(h1) || !finite(h2) || h2 <= 0 || h1 <= h2 || p <= 0 {
		return 0, 0, fmt.Errorf("integration: Рунге-Ромберг требует конечных значений, h1 > h2 > 0 и p > 0")
	}
	factor := math.Pow(h1/h2, float64(p)) - 1
	if !finite(factor) || factor <= 0 {
		return 0, 0, fmt.Errorf("integration: некорректный коэффициент Рунге-Ромберга")
	}
	err2 = (j2 - j1) / factor
	refined = j2 + err2
	if !finite(refined) || !finite(err2) {
		return 0, 0, fmt.Errorf("integration: неконечное уточнение Рунге-Ромберга")
	}
	return refined, err2, nil
}

func finite(v float64) bool {
	return !math.IsNaN(v) && !math.IsInf(v, 0)
}
