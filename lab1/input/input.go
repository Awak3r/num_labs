package input

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"lab1/matrix"
)

func numbers(path string) ([]float64, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("input: %w", err)
	}
	defer f.Close()
	var nums []float64
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		for _, tok := range strings.Fields(sc.Text()) {
			v, err := strconv.ParseFloat(tok, 64)
			if err != nil {
				return nil, fmt.Errorf("input: %s: '%s': %w", path, tok, err)
			}
			nums = append(nums, v)
		}
	}
	if err := sc.Err(); err != nil {
		return nil, fmt.Errorf("input: %s: %w", path, err)
	}
	return nums, nil
}

func read(path string, withB, withEps bool) (*matrix.Matrix, *matrix.Matrix, float64, error) {
	nums, err := numbers(path)
	if err != nil {
		return nil, nil, 0, err
	}
	n, err := dimension(nums, path)
	if err != nil {
		return nil, nil, 0, err
	}
	need := 1 + n*n
	if withB {
		need += n
	}
	if withEps {
		need++
	}
	if len(nums) < need {
		return nil, nil, 0, fmt.Errorf("input: %s: нужно %d чисел (n, n² элементов A, b если есть, eps если есть), найдено %d", path, need, len(nums))
	}
	a := matrix.Zeros(n, n)
	pos := 1
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			a.Set(i, j, nums[pos])
			pos++
		}
	}
	var b *matrix.Matrix
	if withB {
		b = matrix.Zeros(n, 1)
		for i := 0; i < n; i++ {
			b.Set(i, 0, nums[pos])
			pos++
		}
	}
	var eps float64
	if withEps {
		eps = nums[pos]
	}
	return a, b, eps, nil
}

func ReadMatrix(path string) (*matrix.Matrix, error) {
	a, _, _, err := read(path, false, false)
	return a, err
}

func ReadSystem(path string) (*matrix.Matrix, *matrix.Matrix, error) {
	a, b, _, err := read(path, true, false)
	return a, b, err
}

func ReadMatrixEps(path string) (*matrix.Matrix, float64, error) {
	a, _, eps, err := read(path, false, true)
	return a, eps, err
}

func ReadSystemEps(path string) (*matrix.Matrix, *matrix.Matrix, float64, error) {
	a, b, eps, err := read(path, true, true)
	return a, b, eps, err
}

func ReadTridiagonal(path string) ([]float64, []float64, []float64, []float64, error) {
	nums, err := numbers(path)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	n, err := dimension(nums, path)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	need := 1 + n + (n - 1) + (n - 1) + n
	if len(nums) < need {
		return nil, nil, nil, nil, fmt.Errorf("input: %s: нужно %d чисел (n, главная, верхняя, нижняя, правая часть), найдено %d", path, need, len(nums))
	}
	main := append([]float64(nil), nums[1:1+n]...)
	super := append([]float64(nil), nums[1+n:2*n]...)
	sub := append([]float64(nil), nums[2*n:3*n-1]...)
	d := append([]float64(nil), nums[3*n-1:need]...)
	return main, super, sub, d, nil
}

func dimension(nums []float64, path string) (int, error) {
	if len(nums) == 0 {
		return 0, fmt.Errorf("input: %s: пустой файл", path)
	}
	n := int(nums[0])
	if float64(n) != nums[0] || n <= 0 {
		return 0, fmt.Errorf("input: %s: первое число должно быть целой размерностью, найдено %v", path, nums[0])
	}
	return n, nil
}
