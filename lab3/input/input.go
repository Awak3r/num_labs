package input

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Points struct {
	X    []float64
	Y    []float64
	XStar float64
}

func ReadPointsXstar(path string) (Points, error) {
	nums, err := Numbers(path)
	if err != nil {
		return Points{}, err
	}
	if len(nums) < 2 {
		return Points{}, fmt.Errorf("input: %s: нужно n, x[], y[], x*, найдено %d чисел", path, len(nums))
	}
	n := int(nums[0])
	if float64(n) != nums[0] || n <= 0 {
		return Points{}, fmt.Errorf("input: %s: первое число должно быть целым n, найдено %v", path, nums[0])
	}
	need := 1 + 2*n + 1
	if len(nums) < need {
		return Points{}, fmt.Errorf("input: %s: нужно %d чисел (n, x[], y[], x*), найдено %d", path, need, len(nums))
	}
	return Points{
		X:    nums[1 : 1+n],
		Y:    nums[1+n : 1+2*n],
		XStar: nums[1+2*n],
	}, nil
}

// ReadNodesXstar читает n, x[] (n штук), x* — для задач, где y вычисляется из функции.
func ReadNodesXstar(path string) (x []float64, xstar float64, err error) {
	nums, err := Numbers(path)
	if err != nil {
		return nil, 0, err
	}
	if len(nums) < 1 {
		return nil, 0, fmt.Errorf("input: %s: пустой файл", path)
	}
	n := int(nums[0])
	if float64(n) != nums[0] || n <= 0 {
		return nil, 0, fmt.Errorf("input: %s: первое число должно быть целым n, найдено %v", path, nums[0])
	}
	need := 1 + n + 1
	if len(nums) < need {
		return nil, 0, fmt.Errorf("input: %s: нужно %d чисел (n, x[], x*), найдено %d", path, need, len(nums))
	}
	return nums[1 : 1+n], nums[1+n], nil
}

func ReadPoints(path string) (Points, error) {
	nums, err := Numbers(path)
	if err != nil {
		return Points{}, err
	}
	if len(nums) < 1 {
		return Points{}, fmt.Errorf("input: %s: пустой файл", path)
	}
	n := int(nums[0])
	if float64(n) != nums[0] || n <= 0 {
		return Points{}, fmt.Errorf("input: %s: первое число должно быть целым n, найдено %v", path, nums[0])
	}
	need := 1 + 2*n
	if len(nums) < need {
		return Points{}, fmt.Errorf("input: %s: нужно %d чисел (n, x[], y[]), найдено %d", path, need, len(nums))
	}
	return Points{
		X: nums[1 : 1+n],
		Y: nums[1+n : 1+2*n],
	}, nil
}

func ReadIntegral(path string) (x0, xk, h1, h2 float64, err error) {
	nums, err := Numbers(path)
	if err != nil {
		return 0, 0, 0, 0, err
	}
	if len(nums) < 4 {
		return 0, 0, 0, 0, fmt.Errorf("input: %s: нужно 4 числа (x0, xk, h1, h2), найдено %d", path, len(nums))
	}
	return nums[0], nums[1], nums[2], nums[3], nil
}

func Numbers(path string) ([]float64, error) {
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
