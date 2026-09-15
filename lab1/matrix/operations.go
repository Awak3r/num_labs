package matrix

import (
	"fmt"
	"math"
	"strings"
)

func (m *Matrix) Rows() int { return len(m.body) }
func (m *Matrix) Cols() int { return len(m.body[0]) }

func (m *Matrix) Get(i, j int) float64    { return m.body[i][j] }
func (m *Matrix) Set(i, j int, v float64) { m.body[i][j] = v }

func (m *Matrix) Add(o *Matrix) (*Matrix, error) {
	if m.Rows() != o.Rows() || m.Cols() != o.Cols() {
		return nil, fmt.Errorf("Add: размеры не совпадают (%dx%d и %dx%d)",
			m.Rows(), m.Cols(), o.Rows(), o.Cols())
	}
	res := Zeros(m.Rows(), m.Cols())
	for i := 0; i < m.Rows(); i++ {
		for j := 0; j < m.Cols(); j++ {
			res.body[i][j] = m.body[i][j] + o.body[i][j]
		}
	}
	return res, nil
}

func (m *Matrix) Multiply(o *Matrix) (*Matrix, error) {
	if m.Cols() != o.Rows() {
		return nil, fmt.Errorf("Multiply: несовместимые размеры %dx%d и %dx%d",
			m.Rows(), m.Cols(), o.Rows(), o.Cols())
	}
	res := Zeros(m.Rows(), o.Cols())
	for i := 0; i < m.Rows(); i++ {
		for j := 0; j < o.Cols(); j++ {
			var s float64
			for k := 0; k < m.Cols(); k++ {
				s += m.body[i][k] * o.body[k][j]
			}
			res.body[i][j] = s
		}
	}
	return res, nil
}

func (m *Matrix) Sub(o *Matrix) (*Matrix, error) {
	if m.Rows() != o.Rows() || m.Cols() != o.Cols() {
		return nil, fmt.Errorf("Sub: размеры не совпадают (%dx%d и %dx%d)",
			m.Rows(), m.Cols(), o.Rows(), o.Cols())
	}
	res := Zeros(m.Rows(), m.Cols())
	for i := 0; i < m.Rows(); i++ {
		for j := 0; j < m.Cols(); j++ {
			res.body[i][j] = m.body[i][j] - o.body[i][j]
		}
	}
	return res, nil
}

func (m *Matrix) Norm2() float64 { //чет гавно
	var s float64
	for i := 0; i < m.Rows(); i++ {
		for j := 0; j < m.Cols(); j++ {
			s += m.body[i][j] * m.body[i][j]
		}
	}
	return math.Sqrt(s)
}

func (m *Matrix) MaxAbs() float64 {
	var max float64
	for i := 0; i < m.Rows(); i++ {
		for j := 0; j < m.Cols(); j++ {
			if v := math.Abs(m.body[i][j]); v > max {
				max = v
			}
		}
	}
	return max
}

func (m *Matrix) Transpose() *Matrix {
	res := Zeros(m.Cols(), m.Rows())
	for i := 0; i < m.Rows(); i++ {
		for j := 0; j < m.Cols(); j++ {
			res.body[j][i] = m.body[i][j]
		}
	}
	return res
}

func (m *Matrix) Copy() *Matrix { return Create(m.body) }

func (m *Matrix) SwapRows(i, j int) error {
	if i < 0 || i >= m.Rows() || j < 0 || j >= m.Rows() {
		return fmt.Errorf("SwapRows: индекс вне диапазона (i=%d, j=%d), строк: %d", i, j, m.Rows())
	}
	m.body[i], m.body[j] = m.body[j], m.body[i]
	return nil
}

func (m *Matrix) Print() { fmt.Print(m.String()) }

func (m *Matrix) String() string {
	var b strings.Builder
	for _, row := range m.body {
		fmt.Fprintf(&b, "%v\n", row)
	}
	return b.String()
}
