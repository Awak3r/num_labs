package integration

import "fmt"

// Method — численный метод интегрирования.
type Method struct {
	Name string
	P    int // порядок точности для Рунге–Ромберга
	Fn   func(f func(float64) float64, a, b float64, n int) (float64, error)
}

// Methods — все требуемые методы.
var Methods = []Method{
	{"левые прямоугольники", 1, LeftRect},
	{"правые прямоугольники", 1, RightRect},
	{"средние прямоугольники", 2, MidRect},
	{"трапеции", 2, Trapezoid},
	{"Симпсон", 4, Simpson},
}

// LeftRect — метод левых прямоугольников.
func LeftRect(f func(float64) float64, a, b float64, n int) (float64, error) {
	return rects(f, a, b, n, 0)
}

// RightRect — метод правых прямоугольников.
func RightRect(f func(float64) float64, a, b float64, n int) (float64, error) {
	return rects(f, a, b, n, 1)
}

// MidRect — метод средних прямоугольников.
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

// Trapezoid — метод трапеций.
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

// Simpson — метод Симпсона (парабол), требует чётного n.
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

// RungeRomberg уточняет значение J2 (мелкий шаг) по J1 (крупный шаг):
//
//	J_RR = J2 + (J2 − J1) / (2^p − 1)
//	R — оценка погрешности J2: |J_RR − J2| / (2^p − 1) … фактически |J2−J1|/(2^p−1).
func RungeRomberg(j1, j2 float64, p int) (refined, err2 float64) {
	factor := 1
	for i := 0; i < p; i++ {
		factor *= 2
	}
	factor--
	refined = j2 + (j2-j1)/float64(factor)
	err2 = (j2 - j1) / float64(factor)
	return refined, err2
}
