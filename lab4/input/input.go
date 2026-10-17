package input

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// ReadInterval читает x0, xk, h1, h2 — отрезок и два шага сетки
// (h2 — для оценки погрешности по методу Рунге-Ромберга).
func ReadInterval(path string) (x0, xk, h1, h2 float64, err error) {
	nums, err := Numbers(path)
	if err != nil {
		return 0, 0, 0, 0, err
	}
	if len(nums) < 4 {
		return 0, 0, 0, 0, fmt.Errorf("input: %s: нужно 4 числа (x0, xk, h1, h2), найдено %d", path, len(nums))
	}
	return nums[0], nums[1], nums[2], nums[3], nil
}

// ReadBoundary читает x0, xk, h1, h2, eps — отрезок, два шага сетки
// и точность метода секущих в стрельбе.
func ReadBoundary(path string) (x0, xk, h1, h2, eps float64, err error) {
	nums, err := Numbers(path)
	if err != nil {
		return 0, 0, 0, 0, 0, err
	}
	if len(nums) < 5 {
		return 0, 0, 0, 0, 0, fmt.Errorf("input: %s: нужно 5 чисел (x0, xk, h1, h2, eps), найдено %d", path, len(nums))
	}
	return nums[0], nums[1], nums[2], nums[3], nums[4], nil
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
