package matrix

// Eps — порог «математического нуля» для численных методов.
const Eps = 1e-12

type Matrix struct {
	body [][]float64
}

func Create(m [][]float64) *Matrix {
	body := make([][]float64, len(m))
	for i := range m {
		body[i] = make([]float64, len(m[i]))
		copy(body[i], m[i])
	}
	return &Matrix{body}
}

func Column(v []float64) *Matrix {
	m := Zeros(len(v), 1)
	for i, x := range v {
		m.body[i][0] = x
	}
	return m
}

func Diagonal(v []float64) *Matrix {
	m := Zeros(len(v), len(v))
	for i, x := range v {
		m.body[i][i] = x
	}
	return m
}

func Zeros(rows, cols int) *Matrix {
	body := make([][]float64, rows)
	for i := range body {
		body[i] = make([]float64, cols)
	}
	return &Matrix{body}
}

func MakeIdentityMatrix(order int) *Matrix {
	m := Zeros(order, order)
	for i := 0; i < order; i++ {
		m.body[i][i] = 1
	}
	return m
}
